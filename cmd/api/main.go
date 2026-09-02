package main

import (
	"blessdarah/tuts/internal/app"
	"blessdarah/tuts/internal/config"
	"log"
)

func main() {

	cfg := config.LoadConfig()

	a, err := app.NewApp(cfg)

	if err != nil {
		log.Fatal(err)
	}

	a.Run()

}
