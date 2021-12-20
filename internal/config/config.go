package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// To avoiding problems, config is imported by other parts of the application, but it doesn't import anything else from the application self.
// Config only uses the standard library.

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger // write log
	ErrorLog      *log.Logger
	InProduction  bool // in production or development
	Session       *scs.SessionManager
}
