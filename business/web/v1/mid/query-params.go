package mid

import (
	"context"
	"fmt"
	"net/http"

	v1 "github.com/thebigyovadiaz/golang-base-structure/business/web/v1"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/web"
)

const (
	APIKEY string = "apikey"
	STATUS string = "status"
	LIMIT  string = "limit"
	KeyId  string = "key_id"
	YEAR   string = "year"
	MONTH  string = "month"
)

var queryParamsDict = map[string]string{
	APIKEY: APIKEY,
	STATUS: STATUS,
	LIMIT:  LIMIT,
	KeyId:  KeyId,
	YEAR:   YEAR,
	MONTH:  MONTH,
}

// QueryParams is a middleware that check and valid query params keys.
func QueryParams() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var msg string

			// Checking query params
			for key, _ := range r.URL.Query() {
				_, ok := queryParamsDict[key]
				if !ok {
					msg = fmt.Sprintf("%s: query param [%s] don't allow", http.StatusText(http.StatusUnauthorized), key)
					return v1.NewRequestError(
						fmt.Errorf(msg),
						http.StatusUnauthorized,
						msg)
				}
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
