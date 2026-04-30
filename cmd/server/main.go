package main

import (
	"bam-bo/internal/app"
	"bam-bo/internal/store"
	"log"
)

func main() {
	projectStore, err := store.NewProjectStore("data/projects")
	if err != nil {
		log.Fatal(err)
	}

	server := app.New(projectStore, app.Config{
		TemplatePath: "web/templates/index.html",
		StaticDir:    "web/static",
	})
	if err := server.Run(":8085"); err != nil {
		log.Fatal(err)
	}
}
