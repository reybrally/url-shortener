package main

import (
	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/reybrally/REST-API-app/internal/config"
	"github.com/reybrally/REST-API-app/internal/http-server/handlers/redirect"
	"github.com/reybrally/REST-API-app/internal/http-server/handlers/url/save"
	mwLogger "github.com/reybrally/REST-API-app/internal/http-server/middleware/logger"
	"github.com/reybrally/REST-API-app/internal/lib/logger/handlers/slogpretty"
	"github.com/reybrally/REST-API-app/internal/lib/logger/sl"
	"github.com/reybrally/REST-API-app/internal/storage/sqlite"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	color.NoColor = false
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting utl-shortener", slog.String("env", cfg.Env))
	log.Error("error messages are enable")
	log.Debug("debug messages are enable")
	storage, err := sqlite.New(cfg.Storage)
	if err != nil {
		log.Error(
			"failed to init message", sl.Err(err),
		)
		os.Exit(1)
	}
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/urls", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage))
	})

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))
	//router.Delete("/url/{alias}", redirect.New(log, storage))
	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("failed to start server")
	}
	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
