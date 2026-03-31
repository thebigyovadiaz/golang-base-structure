package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/thebigyovadiaz/golang-base-structure/foundation/logger"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/web"
)

// Logger is a middleware that logs information about incoming HTTP requests and outgoing responses using the provided logger.
// It wraps the given handler and adds logging functionality to it.
func Logger(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v := web.GetValues(ctx)

			log.Info(ctx, "request started "+r.RequestURI, "protocol", r.Proto, "method", r.Method, "URL", r.RequestURI, "userAgent", r.UserAgent(),
				"remoteAddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Info(ctx, "request completed "+r.RequestURI, "protocol", r.Proto, "method", r.Method, "URL", r.RequestURI, "userAgent", r.UserAgent(),
				"remoteAddr", r.RemoteAddr, "statusCode", v.StatusCode, "since", time.Since(v.Now).String())

			return err
		}

		return h
	}

	return m
}
