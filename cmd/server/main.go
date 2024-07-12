package main

import (
	"log"

	"github.com/dub-otrezkov/test_go/internal/api"
	"github.com/dub-otrezkov/test_go/internal/app"
	db "github.com/dub-otrezkov/test_go/internal/database"
	"github.com/dub-otrezkov/test_go/internal/tasks"
)

func main() {

	db, err := db.New()

	if err != nil {
		panic(err)
	}

	API := api.New(db)
	tasksapp := tasks.New()

	port := ":52"
	a := app.New(port, API, tasksapp)

	log.Printf("hosted on  localhost%v\n", port)
	a.Run()
}
