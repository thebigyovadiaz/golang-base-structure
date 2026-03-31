package basestructuregrp

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/thebigyovadiaz/golang-base-structure/business/web/v1/mid"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/logger"

	"github.com/thebigyovadiaz/golang-base-structure/foundation/web"
)

type Config struct {
	Log *logger.Logger
	DB  *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, group string, cfg Config) {
	baseStructureHandler := New("success")

	app.Handle(
		http.MethodGet,
		group,
		"/new-message",
		baseStructureHandler.NewMessage,
		mid.Context(),
	)
}
