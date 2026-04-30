package handlers

import (
	"bam-bo/internal/models"
	"bam-bo/internal/pdf"
	"encoding/json"
	"net/http"
)

func (h *Handler) ExportJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req models.ExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	req.Materials = models.NormalizeMaterials(req.Materials)
	req.Textures = models.NormalizeTextures(req.Textures)

	pieces, err := pdf.GatherPieces(req.Cabinets, req.Materials)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Disposition", `attachment; filename="cut-list.json"`)
	writeJSON(w, http.StatusOK, models.ExportResponse{
		ProjectName: req.ProjectName,
		Cabinets:    req.Cabinets,
		Pieces:      pieces,
		Materials:   req.Materials,
		Textures:    req.Textures,
	})
}

func (h *Handler) ExportPDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req models.ExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	req.Materials = models.NormalizeMaterials(req.Materials)
	req.Textures = models.NormalizeTextures(req.Textures)

	pieces, err := pdf.GatherPieces(req.Cabinets, req.Materials)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	document, err := pdf.Build(req.ProjectName, req.Cabinets, pieces, req.Materials)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to build PDF"})
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="cut-list.pdf"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(document)
}
