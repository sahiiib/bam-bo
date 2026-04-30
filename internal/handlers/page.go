package handlers

import (
	"net/http"
	"os"
)

func (h *Handler) Index(w http.ResponseWriter, _ *http.Request) {
	page, err := os.ReadFile(h.templatePath)
	if err != nil {
		http.Error(w, "failed to load page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(page)
}
