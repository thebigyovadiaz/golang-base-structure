package web

import (
	"context"
	//"gitlab.com/ccla/tapp/kids-banking/go-ms-quests/business/core/user"
	"time"
)

type contextKey int

const ctxKey contextKey = 1
const defaultTraceId = "00000000-0000-000000000000"

// Values struct represents the state for each request.
type Values struct {
	TraceID        string
	CorrelationID  string
	Now            time.Time
	StatusCode     int
	Response       any
	RUT            string
	SecurityToken  string
	Token          string
	RequestChannel string
	WalletToken    string
	WalletRUT      string
	//QueryFilters   user.FilterQueryDB
}

/*
	The following methods are an exception to the policy of no getters & setters.
*/

// GetValues returns the values from the context.
func GetValues(ctx context.Context) *Values {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return &Values{
			TraceID: defaultTraceId,
			Now:     time.Now().UTC(),
		}
	}

	return v
}

// GetTraceID returns the trace ID from the context.
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return defaultTraceId
	}

	return v.TraceID
}

// GetSecurityToken returns the x-security-token from the context.
func GetSecurityToken(ctx context.Context) string {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return ""
	}

	return v.SecurityToken
}

// GetCorrelationID returns the correlation ID from the context.
func GetCorrelationID(ctx context.Context) string {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return defaultTraceId
	}

	return v.CorrelationID
}

// GetTime returns the time from the context.
func GetTime(ctx context.Context) time.Time {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return time.Now()
	}

	return v.Now
}

// SetStatusCode sets the status code back into the context.
func SetStatusCode(ctx context.Context, statusCode int) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.StatusCode = statusCode
}

// SetStatusCode sets the status code back into the context.
func SetResponse(ctx context.Context, response any) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.Response = response
}

// SetRut sets the user's RUT into the context.
func SetRut(ctx context.Context, rut string) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.RUT = rut
}

// GetRut returns the JWT token from the context.
func GetRut(ctx context.Context) string {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return ""
	}

	return v.RUT
}

func SetRequestChannel(ctx context.Context, reqChannel string) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}
	v.RequestChannel = reqChannel

}

// SetSecurityToken assigns the user's security token to the context.
// Storing this value in the context ensures the availability of a security token,
// even in case of service failure. This prevents automatic user logout by the mobile
// application following the invocation of this service.
func SetSecurityToken(ctx context.Context, token string) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.SecurityToken = token
}

// SetToken sets the user's JWT Token into the context.
func SetToken(ctx context.Context, token string) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.Token = token
}

// GetToken returns the JWT token from the context.
func GetToken(ctx context.Context) string {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return ""
	}

	return v.Token
}

// SetWalletRut sets the user's RUT into the context.
func SetWalletRut(ctx context.Context, rut string) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.WalletRUT = rut
}

// GetWalletRut returns the JWT token from the context.
func GetWalletRut(ctx context.Context) string {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return ""
	}

	return v.WalletRUT
}

// SetWalletToken sets the user's JWT Token into the context.
func SetWalletToken(ctx context.Context, token string) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.WalletToken = token
}

// GetWalletToken returns the JWT token from the context.
func GetWalletToken(ctx context.Context) string {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return ""
	}

	return v.WalletToken
}

// SetQueryFilters sets filters into the context.
/*func SetQueryFilters(ctx context.Context, filters user.FilterQueryDB) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.QueryFilters = filters
}

// GetQueryFilters returns query filters for DB query.
func GetQueryFilters(ctx context.Context) user.FilterQueryDB {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return user.FilterQueryDB{}
	}

	return v.QueryFilters
}*/
