package main

import (
	"github.com/dub-otrezkov/OschApp/internal/api"
	"github.com/dub-otrezkov/OschApp/internal/app"
	db "github.com/dub-otrezkov/OschApp/internal/database"
	"github.com/dub-otrezkov/OschApp/internal/tasks"
	"github.com/dub-otrezkov/OschApp/internal/userstats"
	"github.com/dub-otrezkov/OschApp/pkg/auth"
)

func main() {

	db, err := db.New()

	if err != nil {
		panic(err)
	}

	API_app := api.New(db)
	auth_app := auth.New(db)
	tasks_app := tasks.New(db)
	userstats_app := userstats.New()

	port := ":52"
	a := app.New(port, API_app, auth_app, tasks_app, userstats_app)

	a.Run()
}
