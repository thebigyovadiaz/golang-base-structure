package healthgrp

import (
	"net/http"
	"time"

	"github.com/thebigyovadiaz/golang-base-structure/foundation/logger"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/web"
)

type Config struct {
	Since time.Time
	Build string
	Cores string
	Log   *logger.Logger
}

// Routes adds specific routes for this group.
func Routes(app *web.App, group string, cfg Config) {

	hgh := New(cfg.Build, cfg.Since, cfg.Cores)

	app.CustomHandle(http.MethodGet, group, "/health", hgh.Health)
}
