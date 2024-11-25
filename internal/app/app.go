package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kleo-53/music-system/config"
	v1 "github.com/kleo-53/music-system/internal/controller"
	"github.com/kleo-53/music-system/internal/migrate"
	songService "github.com/kleo-53/music-system/internal/service/song"
	songStore "github.com/kleo-53/music-system/internal/store/song"
	"github.com/kleo-53/music-system/pkg/logger"
	"github.com/kleo-53/music-system/pkg/postgres"
)

func Run(cfg *config.Config) {
	if err := godotenv.Load("config.env"); err != nil {
		logger.Log().Warn(context.Background(), "No .env file found, using environment variables")
	}

	ctx := context.Background()
	logger.New(cfg.LogLevel)
	err := migrate.CreateDBIfNotExists(ctx, os.Getenv("DB_ADMIN_URL"), os.Getenv("DB_NAME"))
	if err != nil {
		logger.Log().Fatal(ctx, "failed to create database: %s", err.Error())
		return
	}
	pg, err := postgres.New(ctx, cfg.DBURL)
	if err != nil {
		logger.Log().Fatal(ctx, "error with connection to database: %s", err.Error())
	}
	defer pg.Close(ctx)

	if err := migrate.RunMigration(cfg.DBURL); err != nil {
		logger.Log().Fatal(ctx, "error with up migrations for database: %s", err.Error())
		return
	}
	songStore := songStore.New(pg)
	songService := songService.New(songStore)

	app := mux.NewRouter()
	v1.NewRouter(
		app,
		songService,
	)
	server := &http.Server{
		Addr: cfg.Port,
	}
	logger.Log().Info(ctx, "server was started on %s", cfg.Port)
	http.Handle("/", app)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Log().Fatal(ctx, "HTTP server error: %v", err)
		}
		logger.Log().Info(ctx, "Stopped serving new connections.")
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log().Fatal(shutdownCtx, "HTTP shutdown error: %v", err)
	}
	logger.Log().Info(ctx, "Graceful shutdown complete.")
}
