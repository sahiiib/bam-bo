package store

import (
	"bam-bo/internal/models"
	"bam-bo/internal/pdf"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type ProjectStore struct {
	dir string
}

func NewProjectStore(dir string) (*ProjectStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	return &ProjectStore{dir: dir}, nil
}

func (s *ProjectStore) Save(payload models.ProjectPayload) (models.Project, error) {
	now := time.Now().UTC()
	projectID := sanitizeProjectID(payload.ID)
	if projectID == "" {
		projectID = sanitizeProjectID(payload.Name)
	}
	if projectID == "" {
		projectID = "project-" + now.Format("20060102-150405")
	}

	project := models.Project{
		ID:        projectID,
		Name:      strings.TrimSpace(payload.Name),
		Cabinets:  payload.Cabinets,
		Materials: models.NormalizeMaterials(payload.Materials),
		Textures:  models.NormalizeTextures(payload.Textures),
		UpdatedAt: now,
	}

	if project.Name == "" {
		project.Name = "Untitled Project"
	}

	existing, err := s.Get(projectID)
	if err == nil {
		project.CreatedAt = existing.CreatedAt
	} else {
		project.CreatedAt = now
	}

	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return models.Project{}, err
	}

	if err := os.WriteFile(s.filePath(projectID), data, 0o644); err != nil {
		return models.Project{}, err
	}

	return project, nil
}

func (s *ProjectStore) Get(id string) (models.Project, error) {
	data, err := os.ReadFile(s.filePath(id))
	if err != nil {
		return models.Project{}, err
	}

	var project models.Project
	if err := json.Unmarshal(data, &project); err != nil {
		return models.Project{}, err
	}

	project.Materials = models.NormalizeMaterials(project.Materials)
	project.Textures = models.NormalizeTextures(project.Textures)

	return project, nil
}

func (s *ProjectStore) List() ([]models.ProjectSummary, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	summaries := make([]models.ProjectSummary, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		project, err := s.Get(strings.TrimSuffix(entry.Name(), ".json"))
		if err != nil {
			continue
		}

		pieces, err := pdf.GatherPieces(project.Cabinets, project.Materials)
		if err != nil {
			continue
		}

		summaries = append(summaries, models.ProjectSummary{
			ID:           project.ID,
			Name:         project.Name,
			CabinetCount: len(project.Cabinets),
			PieceCount:   pdf.TotalPieceCount(pieces),
			UpdatedAt:    project.UpdatedAt,
		})
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].UpdatedAt.After(summaries[j].UpdatedAt)
	})

	return summaries, nil
}

func (s *ProjectStore) filePath(id string) string {
	return filepath.Join(s.dir, sanitizeProjectID(id)+".json")
}

func sanitizeProjectID(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}

	var b strings.Builder
	lastDash := false
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		case r == '-' || r == '_' || r == ' ':
			if !lastDash && b.Len() > 0 {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}

	out := strings.Trim(b.String(), "-")
	return out
}

func IsNotFound(err error) bool {
	return errors.Is(err, os.ErrNotExist)
}
