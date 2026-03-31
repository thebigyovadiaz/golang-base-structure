package mid

import (
	"context"
	"net/http"

	v1 "github.com/thebigyovadiaz/golang-base-structure/business/web/v1"
	"github.com/thebigyovadiaz/golang-base-structure/business/web/v1/response"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/logger"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/modelvalidator"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/web"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				var er response.Error
				var status int

				v := web.GetValues(ctx)
				securityToken := r.Header.Get("x-security-token")

				switch {
				// ERROR 400 - Fields error
				case modelvalidator.IsFieldErrors(err):
					fieldErrors := modelvalidator.GetFieldErrors(err)

					status = http.StatusBadRequest
					er.Message = http.StatusText(status)
					er.TraceID = v.TraceID
					er.SecurityToken = securityToken
					er.Success = false
					er.AdditionalInfo = fieldErrors.Fields()

				// ERROR 400 - Request error
				case v1.IsRequestError(err):
					reqErr := v1.GetRequestError(err)
					status = reqErr.Status
					er.Message = reqErr.CustomMessage
					er.TraceID = v.TraceID
					er.SecurityToken = securityToken
					er.Success = false

				// ERROR 404 - Not found
				/*case userdb.IsUserNotFoundError(err):
				status = http.StatusNotFound
				er.Message = "User not found"
				er.TraceID = v.TraceID
				er.SecurityToken = securityToken
				er.Success = false*/

				// ERROR 500 - Internal Server error
				default:
					status = http.StatusInternalServerError
					er.Message = http.StatusText(status)
					er.TraceID = v.TraceID
					er.SecurityToken = securityToken
					er.Success = false

				}

				log.Errorc(ctx, 4, "handling error coming out of the call chain",
					"code", status, "errMessage", er.Message, "errDetails", err.Error())

				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
