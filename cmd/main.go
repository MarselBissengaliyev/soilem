package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/MarselBissengaliyev/soilem/configs"
	"github.com/MarselBissengaliyev/soilem/internal/handlers"
	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/MarselBissengaliyev/soilem/internal/repo"
	"github.com/MarselBissengaliyev/soilem/internal/service"
	log "github.com/sirupsen/logrus"
)

// @title Soilem API
// @version 1.0
// @description Api Server for Soilem Application

// @contact.name Marsel Bissengaliyev
// @contact.url https://t.me/marsel_bissengaliyev
// @contact.email marselbisengaliev1@gmail.com

// @host localhost:8000
// BasePath /
func main() {
	// Creating context
	ctx, cancel := context.WithCancel(context.Background())
	// Cancel context at the end
	defer cancel()

	// Set JSON logger formatter
	log.SetFormatter(new(log.JSONFormatter))

	// Init config
	cfg, err := configs.InitConfig("./configs", "yaml", "config")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create connection to postgres DB
	db, err := repo.NewPostgresDB(&repo.PostgresConfig{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		DBName:   cfg.Postgres.DBName,
		UserName: cfg.Postgres.UserName,
		Password: cfg.Postgres.Password,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	repos := repo.NewRepository(db)
	services := service.NewService(repos, cfg)
	handler := handlers.NewHandler(services)

	srv := new(model.Server)
	// Run server in goroutine and init routes
	go func() {
		if err := srv.Run(cfg.Port, handler.InitRoutes()); err != nil {
			log.Fatal("Failed to run http server: ", err)
			return
		}
	}()

	log.Print("Soilem App Started")

	// Wait signals SITERM and SIGINT
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Soilem App Shutting down")

	// If signal came make shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Failed to shutdown http server: ", err)
	}

	// Close db connection
	if err := db.Close(ctx); err != nil {
		log.Error("Failed to to close db connection")
	}
}
