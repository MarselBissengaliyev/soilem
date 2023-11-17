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
	db, err := repo.NewPostgresDB(&repo.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		DBName:   cfg.Database.DBName,
		UserName: cfg.Database.Username,
		Password: cfg.Database.Password,
	}, ctx)

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
