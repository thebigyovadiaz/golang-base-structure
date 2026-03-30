// Package healthgrp holds a handler for health checking.
package healthgrp

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/timecl"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/web"
)

const (
	UP = "UP"
)

// Handlers manages the health endpoint used by k8s.
type Handlers struct {
	build string
	since time.Time
	cores string
	db    *sqlx.DB
}

// New constructs a Handlers type for route access.
func New(build string, since time.Time, cores string) *Handlers {
	return &Handlers{
		build: build,
		since: since,
		cores: cores,
	}
}

// Health allows k8s to know if the service is running.
//
//	@Summary		Health check
//	@Tags			account
//	@Accept			json
//	@Produce		json
//
//	@Success		200	{object} HealthResponse
//
//	@Failure		500	{object} any
//	@Router			/health [get]
func (h *Handlers) Health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	now := timecl.Now()
	uptime := now.Sub(h.since).Truncate(time.Second)
	uptimestr := uptime.String()

	data := HealthResponse{
		Status:     UP,
		Uptime:     uptimestr,
		Timestamp:  now,
		Since:      h.since,
		Version:    h.build,
		GOMAXPROCS: h.cores,
	}

	return web.Respond(ctx, w, data, http.StatusOK)
}
