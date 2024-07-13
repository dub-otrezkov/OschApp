package main

import (
	"github.com/dub-otrezkov/test_go/internal/api"
	"github.com/dub-otrezkov/test_go/internal/app"
	db "github.com/dub-otrezkov/test_go/internal/database"
	"github.com/dub-otrezkov/test_go/internal/tasks"
	"github.com/dub-otrezkov/test_go/pkg/auth"
)

func main() {

	db, err := db.New()

	if err != nil {
		panic(err)
	}

	API := api.New(db)
	tasksapp := tasks.New()
	auth := auth.New()

	port := ":52"
	a := app.New(port, API, tasksapp, auth)

	a.Run()
}
