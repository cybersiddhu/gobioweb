package gobioweb

import (
	dbi "database/sql"
	"github.com/gorilla/sessions"
)

type App struct {
	Template *TemplateWithLayout
	Session  sessions.Store
	Database *dbi.DB
}
