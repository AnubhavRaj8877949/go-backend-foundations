package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"example.com/go-backend-foundation/internal/bootstrap"
)

func main() {
	app, err := bootstrap.New()
	if err != nil {
		log.Fatalf("application initialization failed: %v", err)
	}
		defer app.Close()

	app.Logger.Info(
		"Application initialized",
		"app", app.Config.App.Name,
		"environment", app.Config.App.Environment,
		"port", app.Config.HTTP.Port,
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,    //Terminal mein Ctrl + C
		syscall.SIGTERM, //Docker/Kubernetes application ko stop karte hain
	)

	defer stop()

	if err := app.Run(ctx); err != nil {
		app.Logger.Error("application stopped with error", "error", err)
		os.Exit(1)
	}

}
