package v1

import (
	"github.com/jmoiron/sqlx"
	"github.com/thebigyovadiaz/golang-base-structure/app/api/handlers/v1/basestructuregrp"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/logger"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/web"
)

// Constant values related to the group and version information of the API.
const (
	v1 = "/v1"
)

// Config contains the mandatory values required by the handlers.
type Config struct {
	Log *logger.Logger
	DB  *sqlx.DB
}

// Routes binds all the current version routes.
func Routes(app *web.App, group string, cfg Config) {
	basestructuregrp.Routes(app, group+v1, basestructuregrp.Config{
		Log: cfg.Log,
		DB:  cfg.DB,
	})
}
