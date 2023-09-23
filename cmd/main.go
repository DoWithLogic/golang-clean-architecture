package main

import (
	"context"
	"os"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/internal/app"
)

func main() {
	env := os.Getenv("env")
	if env == "" {
		env = "local"
	}

	cfg, err := config.LoadConfig(env)
	if err != nil {
		panic(err)
	}

	app := app.NewApp(context.Background(), cfg)
	if err := app.Start(); err != nil {
		panic(err)
	}
}
