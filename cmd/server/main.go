package main

import (
	"log/slog"
	"ne-pridumal/effective-mobile-test/internal/api"
	"ne-pridumal/effective-mobile-test/internal/config"
	"ne-pridumal/effective-mobile-test/internal/httpServer"
	pg "ne-pridumal/effective-mobile-test/internal/storage/postgres"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"ne-pridumal/effective-mobile-test/docs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	sl "ne-pridumal/effective-mobile-test/lib/logger/slog"
)

func main() {

	docs.SwaggerInfo.Title = "Test API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http"}

	cfg := config.MustLoad()
	log := sl.New(cfg.Env)

	log.Info("staring restAPI backend", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	cl := api.New(cfg.Api, log)

	db, err := pg.New(cfg.Postgres, log)

	if err != nil {
		log.Error("cannot open postgress connection", sl.Err(err))
	}
	defer db.Close()

	d, err := postgres.WithInstance(db.Db, &postgres.Config{})

	if err != nil {
		log.Error("cannot get postgres driver to migrations", sl.Err(err))
	}

	mg, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", d)
	if err != nil {
		log.Error("error with migrator", sl.Err(err))
	}

	if err := mg.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error("error during migration", sl.Err(err))
	}

	log.Info("setting up server")

	s := httpServer.New(db.Users(), db.Tasks(), cl, log)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	srv := &http.Server{
		Addr:         cfg.HttpServer.Address + ":" + cfg.HttpServer.Port,
		Handler:      s,
		ReadTimeout:  cfg.HttpServer.IdleTimeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", sl.Err(err))
		}
	}()

	log.Info("server started")

	<-done

	log.Info("stopping server")

	log.Info("server stopped")
}
