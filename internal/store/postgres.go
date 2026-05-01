package store

import (
	"bam-bo/internal/models"
	"bam-bo/internal/pdf"
	"context"
	"encoding/json"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresProjectStore struct {
	pool *pgxpool.Pool
}

func NewPostgresProjectStore(ctx context.Context, databaseURL string) (*PostgresProjectStore, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	store := &PostgresProjectStore{pool: pool}
	if err := store.migrate(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return store, nil
}

func (s *PostgresProjectStore) Close() {
	s.pool.Close()
}

func (s *PostgresProjectStore) Save(payload models.ProjectPayload) (models.Project, error) {
	ctx := context.Background()
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

	data, err := json.Marshal(project)
	if err != nil {
		return models.Project{}, err
	}

	_, err = s.pool.Exec(ctx, `
		insert into projects (id, name, document, created_at, updated_at)
		values ($1, $2, $3, $4, $5)
		on conflict (id) do update set
			name = excluded.name,
			document = excluded.document,
			updated_at = excluded.updated_at
	`, project.ID, project.Name, data, project.CreatedAt, project.UpdatedAt)
	if err != nil {
		return models.Project{}, err
	}

	return project, nil
}

func (s *PostgresProjectStore) Get(id string) (models.Project, error) {
	ctx := context.Background()
	var data []byte
	if err := s.pool.QueryRow(ctx, `select document from projects where id = $1`, sanitizeProjectID(id)).Scan(&data); err != nil {
		if err == pgx.ErrNoRows {
			return models.Project{}, os.ErrNotExist
		}
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

func (s *PostgresProjectStore) List() ([]models.ProjectSummary, error) {
	ctx := context.Background()
	rows, err := s.pool.Query(ctx, `select document from projects order by updated_at desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summaries := make([]models.ProjectSummary, 0)
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}

		var project models.Project
		if err := json.Unmarshal(data, &project); err != nil {
			continue
		}
		project.Materials = models.NormalizeMaterials(project.Materials)

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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].UpdatedAt.After(summaries[j].UpdatedAt)
	})

	return summaries, nil
}

func (s *PostgresProjectStore) migrate(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `
		create table if not exists projects (
			id text primary key,
			name text not null,
			document jsonb not null,
			created_at timestamptz not null,
			updated_at timestamptz not null
		);
		create index if not exists projects_updated_at_idx on projects (updated_at desc);
	`)
	return err
}
