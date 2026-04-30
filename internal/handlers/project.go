package handlers

import (
	"bam-bo/internal/models"
	"bam-bo/internal/pdf"
	"bam-bo/internal/store"
	"encoding/json"
	"net/http"
	"strings"
)

type ProjectStore interface {
	Get(id string) (models.Project, error)
	List() ([]models.ProjectSummary, error)
	Save(payload models.ProjectPayload) (models.Project, error)
}

type Handler struct {
	store        ProjectStore
	templatePath string
}

func New(projectStore ProjectStore, templatePath string) *Handler {
	return &Handler{
		store:        projectStore,
		templatePath: templatePath,
	}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) Materials(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"materials": models.DefaultPlyMaterials(),
		"textures":  models.DefaultTextures(),
	})
}

func (h *Handler) Projects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		summaries, err := h.store.List()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list projects"})
			return
		}
		writeJSON(w, http.StatusOK, map[string][]models.ProjectSummary{"projects": summaries})
	case http.MethodPost:
		var payload models.ProjectPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
			return
		}

		if strings.TrimSpace(payload.Name) == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "project name is required"})
			return
		}

		payload.Materials = models.NormalizeMaterials(payload.Materials)
		payload.Textures = models.NormalizeTextures(payload.Textures)

		if _, err := pdf.GatherPieces(payload.Cabinets, payload.Materials); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		project, err := h.store.Save(payload)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save project"})
			return
		}

		writeJSON(w, http.StatusOK, project)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (h *Handler) ProjectByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	id := pathID(r.URL.Path, "/api/projects/")
	if strings.TrimSpace(id) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "project id is required"})
		return
	}

	project, err := h.store.Get(id)
	if err != nil {
		if store.IsNotFound(err) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "project not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load project"})
		return
	}

	writeJSON(w, http.StatusOK, project)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

func pathID(path, prefix string) string {
	return strings.TrimPrefix(path, prefix)
}
