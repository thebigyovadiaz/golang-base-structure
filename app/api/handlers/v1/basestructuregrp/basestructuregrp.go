package basestructuregrp

import (
	"context"
	"net/http"

	"github.com/thebigyovadiaz/golang-base-structure/foundation/web"
)

type Handlers struct {
	Message string
}

func New(msg string) *Handlers {
	return &Handlers{
		Message: msg,
	}
}

// NewMessage receive a message to transform
func (h *Handlers) NewMessage(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	parent := web.GetRut(ctx)
	var message string

	switch h.Message {
	case "success":
		message = "Successfully response"
	default:
		message = "Other message response"
	}

	responseAPI := ChallengeResponse{
		Message:       "Challenge created successfully",
		TraceID:       web.GetTraceID(ctx),
		SecurityToken: r.Header.Get("x-security-token"),
		Success:       true,
		Data: struct {
			Message string `json:"message"`
			RUT     string `json:"rut"`
		}{
			Message: message,
			RUT:     parent,
		},
	}

	return web.Respond(ctx, w, responseAPI, http.StatusAccepted)
}
