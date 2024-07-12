package main

import (
	"fmt"

	"github.com/dub-otrezkov/test_go/internal/api"
	"github.com/dub-otrezkov/test_go/internal/app"
)

func main() {
	fmt.Println("work")

	API := api.API.New()

	a := app.New(":52")
	a.Run()
}
