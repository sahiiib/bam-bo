package app

import "net/http"

func (a *App) routes() {
	a.mux.HandleFunc("/", a.handler.Index)
	a.mux.HandleFunc("/health", a.handler.Health)
	a.mux.HandleFunc("/api/materials", a.handler.Materials)
	a.mux.HandleFunc("/api/projects", a.handler.Projects)
	a.mux.HandleFunc("/api/projects/", a.handler.ProjectByID)
	a.mux.HandleFunc("/api/export/json", a.handler.ExportJSON)
	a.mux.HandleFunc("/api/export/pdf", a.handler.ExportPDF)
	a.mux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir(a.static)),
		),
	)
}
