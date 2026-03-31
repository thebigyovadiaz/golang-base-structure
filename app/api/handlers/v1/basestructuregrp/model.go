package basestructuregrp

import "time"

type ChallengeResponse struct {
	Message       string `json:"message"`
	TraceID       string `json:"traceID"`
	SecurityToken string `json:"securityToken,omitempty"`
	Success       bool   `json:"success"`
	Data          any    `json:"data,omitempty"`
	Metadata      any    `json:"metadata,omitempty"`
}

type ChallengeRequest struct {
	Description    string `json:"description" validate:"required,min=5"`
	ExpirationDate string `json:"expirationDate" validate:"required"`
	Amount         int    `json:"amount" validate:"required"`
	Category       string `json:"category" validate:"required,min=3"`
}

type ChallengeRequestDB struct {
	Description    string    `json:"description"`
	Amount         int       `json:"amount"`
	CategoryId     int       `json:"categoryId"`
	IsActive       bool      `json:"isActive"`
	RelationId     int       `json:"relationId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	ExpirationDate time.Time `json:"expirationDate"`
}

type ChallengesNotFound struct {
	Challenges any `json:"challenges"`
}
