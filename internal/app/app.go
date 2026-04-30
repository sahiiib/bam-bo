package app

import (
	"bam-bo/internal/handlers"
	"net/http"
)

type App struct {
	handler *handlers.Handler
	mux     *http.ServeMux
	static  string
}

type Config struct {
	TemplatePath string
	StaticDir    string
}

func New(projectStore handlers.ProjectStore, cfg Config) *App {
	if cfg.TemplatePath == "" {
		cfg.TemplatePath = "web/templates/index.html"
	}
	if cfg.StaticDir == "" {
		cfg.StaticDir = "web/static"
	}

	app := &App{
		handler: handlers.New(projectStore, cfg.TemplatePath),
		mux:     http.NewServeMux(),
		static:  cfg.StaticDir,
	}

	app.routes()
	return app
}

func (a *App) Run(addr string) error {
	return http.ListenAndServe(addr, withSecurityHeaders(a.mux))
}
