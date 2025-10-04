package main

import (
	"log"

	"github.com/pliffdax/sparrow-api/internal/app"
)

func main() {
	if err := app.New().Run(); err != nil {
		log.Fatal(err)
	}
}
