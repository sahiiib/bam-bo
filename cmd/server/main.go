package main

import (
	"bam-bo/internal/app"
	"bam-bo/internal/models"
	"bam-bo/internal/store"
	"context"
	"log"
	"os"
)

func main() {
	var projectStore interface {
		Get(string) (models.Project, error)
		List() ([]models.ProjectSummary, error)
		Save(models.ProjectPayload) (models.Project, error)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		postgresStore, err := store.NewPostgresProjectStore(context.Background(), databaseURL)
		if err != nil {
			log.Fatal(err)
		}
		defer postgresStore.Close()
		projectStore = postgresStore
	} else {
		filesystemStore, err := store.NewProjectStore("data/projects")
		if err != nil {
			log.Fatal(err)
		}
		projectStore = filesystemStore
	}

	server := app.New(projectStore, app.Config{
		TemplatePath: "web/templates/index.html",
		StaticDir:    "web/static",
	})
	if err := server.Run(":8085"); err != nil {
		log.Fatal(err)
	}
}
