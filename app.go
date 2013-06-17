package gobioweb

import (
	dbi "database/sql"
	"github.com/gorilla/sessions"
	"html/template"
)

type App struct {
	Template *template.Template
	Session  sessions.Store
	Database *dbi.DB
}
