package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/thebigyovadiaz/golang-base-structure/foundation/web"
)

// Panics recovers from a panic and converts it to an error, which allows
// it to be reported and handled in the Errors middleware.
func Panics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			// Defer a function to recover from a panic after executing the next handler.
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()

					// The named return value "err" contains the new error. This way,
					// we can return the error by bypassing the default behavior that
					// doesn't allow returning values.
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))
				}
			}()

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
