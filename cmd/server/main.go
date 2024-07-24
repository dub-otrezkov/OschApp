package main

import (
	"github.com/dub-otrezkov/OschApp/internal/api"
	"github.com/dub-otrezkov/OschApp/internal/app"
	db "github.com/dub-otrezkov/OschApp/internal/database"
	"github.com/dub-otrezkov/OschApp/pkg/auth"
)

func main() {

	db, err := db.New()

	if err != nil {
		panic(err)
	}

	API := api.New(db)
	auth := auth.New(db, "User")

	port := ":52"
	a := app.New(port, API, auth)

	a.Run()
}
