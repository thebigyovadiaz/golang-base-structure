// Package handlers manages the different versions of the API.
package handlers

import (
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/thebigyovadiaz/golang-base-structure/app/api/handlers/healthgrp"
	//v1 "github.com/thebigyovadiaz/golang-base-structure/app/api/handlers/v1"
	//"github.com/thebigyovadiaz/golang-base-structure/business/web/v1/mid"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/logger"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/timecl"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/web"
)

const (
	group = "/quests"
)

type APIMuxConfig struct {
	Build    string
	Shutdown chan os.Signal
	Log      *logger.Logger
	DB       *sqlx.DB
}

// APIMux constructs a http.Handler that contains the app routes.
func APIMux(cfg APIMuxConfig) http.Handler {
	cores := strconv.Itoa(runtime.GOMAXPROCS(0))
	startTime := timecl.Now()

	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Panics())

	healthgrp.Routes(app, group, healthgrp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		Since: startTime,
		Cores: cores,
	})

	v1.Routes(app, group, v1.Config{
		Log:              cfg.Log,
		DB:               cfg.DB,
		WalletAuth:       cfg.WalletAuth,
		PushNotification: cfg.PushNotification,
	})

	return app
}
