package main

import (
	"blessdarah/tuts/internal/app"
	"blessdarah/tuts/internal/config"
)

func main() {

	cfg := config.LoadConfig()

	a := app.NewApp(cfg)

	a.Run()

}
