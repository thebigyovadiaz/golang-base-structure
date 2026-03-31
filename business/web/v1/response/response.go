package response

// Error is the structure used by the API to respond to the client
// when a failure happens.
type Error struct {
	Message        string `json:"message"`
	TraceID        string `json:"traceID"`
	SecurityToken  string `json:"securityToken,omitempty"`
	Success        bool   `json:"success"`
	AdditionalInfo any    `json:"data,omitempty"`
}
