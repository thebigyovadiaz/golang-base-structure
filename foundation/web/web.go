// Package web contains a small web framework extension.
package web

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles an http request.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint of the app. It configures the context object for each http handler.
type App struct {
	*chi.Mux
	shutdown chan os.Signal
	mw       []Middleware
}

// NewApp returns an App value that handles a set of routes for the app.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	mux := chi.NewMux()

	return &App{
		Mux:      mux,
		shutdown: shutdown,
		mw:       mw,
	}
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// Handle associates a handler function with the specified http method and path.
func (a *App) Handle(method, group, path string, handler Handler, mw ...Middleware) {
	switch {
	case len(mw) > 0:
		handler = wrapMiddleware(mw, handler)
		fallthrough
	case len(a.mw) > 0:
		handler = wrapMiddleware(a.mw, handler)
	}

	a.handle(method, group, path, handler)
}

// CustomHandle is similar to Handle function, but it requires you to specify explicitly
// the desired middlewares. No app middlewares are set by default.
func (a *App) CustomHandle(method, group, path string, handler Handler, mw ...Middleware) {
	if len(mw) > 0 {
		handler = wrapMiddleware(mw, handler)
	}

	a.handle(method, group, path, handler)
}

func (a *App) handle(method, group, path string, handler Handler) {
	h := func(w http.ResponseWriter, r *http.Request) {
		cid := r.Header.Get("correlationID")

		// set trace id and init time for the incoming request.
		v := Values{TraceID: uuid.NewString(), CorrelationID: cid, Now: time.Now().UTC()}
		ctx := context.WithValue(r.Context(), ctxKey, &v)

		if err := handler(ctx, w, r); err != nil {
			if validateShutdown(err) {
				a.SignalShutdown()
				return
			}
		}
	}

	a.Mux.MethodFunc(method, group+path, h)
}

// validateShutdown validates the error for special conditions that do not
// warrant an actual shutdown by the system.
func validateShutdown(err error) bool {

	// Ignore syscall.EPIPE and syscall.ECONNRESET errors which occurs
	// when a write operation happens on the http.ResponseWriter that
	// has simultaneously been disconnected by the client (TCP
	// connections is broken). For instance, when large amounts of
	// data is being written or streamed to the client.
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// https://gosamples.dev/broken-pipe/
	// https://gosamples.dev/connection-reset-by-peer/

	switch {
	case errors.Is(err, syscall.EPIPE):

		// Usually, you get the broken pipe error when you write to the connection after the
		// RST (TCP RST Flag) is sent.
		// The broken pipe is a TCP/IP error occurring when you write to a stream where the
		// other end (the peer) has closed the underlying connection. The first write to the
		// closed connection causes the peer to reply with an RST packet indicating that the
		// connection should be terminated immediately. The second write to the socket that
		// has already received the RST causes the broken pipe error.
		return false

	case errors.Is(err, syscall.ECONNRESET):

		// Usually, you get connection reset by peer error when you read from the
		// connection after the RST (TCP RST Flag) is sent.
		// The connection reset by peer is a TCP/IP error that occurs when the other end (peer)
		// has unexpectedly closed the connection. It happens when you send a packet from your
		// end, but the other end crashes and forcibly closes the connection with the RST
		// packet instead of the TCP FIN, which is used to close a connection under normal
		// circumstances.
		return false
	}

	return true
}
