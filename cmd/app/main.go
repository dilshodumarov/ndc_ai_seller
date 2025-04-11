package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"sugurta/internal/app"
	"sugurta/internal/pkg/config"

	"go.uber.org/zap"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app, err := app.NewApp(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	// run application
	go func() {
		app.Logger.Info("Listen: ", zap.String("address", cfg.Server.Host+cfg.Server.Port))
		if err := app.Run(); err != nil {
			app.Logger.Error("app run", zap.Error(err))
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	// app stops
	app.Logger.Info("api gateway service stops")
	app.Stop()
}
