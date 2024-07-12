package main

import (
	"fmt"

	"github.com/dub-otrezkov/test_go/internal/app"
)

func main() {
	fmt.Println("work")

	// d := db.MySQLdatabase.New()

	a := app.New(":52")
	a.Run()
}
