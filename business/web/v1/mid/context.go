package mid

import (
	"context"
	"net/http"

	"github.com/thebigyovadiaz/golang-base-structure/foundation/web"
)

func Context() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			web.SetSecurityToken(ctx, r.Header.Get("x-security-token"))
			web.SetRut(ctx, r.Header.Get("x-customer-rut"))
			web.SetRequestChannel(ctx, r.Header.Get("x-request-channel"))
			return handler(ctx, w, r)
		}
		return h
	}

	return m
}
