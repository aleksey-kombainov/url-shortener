package main

import (
	"context"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/config"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/http"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/logger"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage"
	"github.com/rs/zerolog"
	nethttp "net/http"
	"os"
	"time"
)

func main() {
	var err error

	loggerInstance := logger.GetLogger()

	options := config.GetOptions(loggerInstance)
	if err != nil {
		loggerInstance.Error().Msg("Can't parse env vars: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var storageInstance storage.ShortcutStorager
	if options.DatabaseDsn != "" {
		storageInstance, err = storage.ShortcutStorageFactory(ctx, &loggerInstance, storage.TypeDB, options.DatabaseDsn)
	} else if options.FileStoragePath != "" {
		storageInstance, err = storage.ShortcutStorageFactory(ctx, &loggerInstance, storage.TypeFile, options.FileStoragePath)
	} else {
		storageInstance, err = storage.ShortcutStorageFactory(ctx, &loggerInstance, storage.TypeMemory, "")
	}
	if err != nil {
		loggerInstance.Error().Msgf("cant get storage: %w", err)
		shutdown(&loggerInstance)
	} else {
		defer func() {
			if err := storageInstance.Close(); err != nil {
				loggerInstance.Error().Msgf("can't close storage: %w")
			}
		}()
	}

	loggerInstance.Info().Msg("Starting server")

	shortcutService := app.NewShortcutService(&loggerInstance, &storageInstance)
	urlManagerService, err := app.NewURLManagerServiceFromFullURL(options.BaseURL)
	if err != nil {
		loggerInstance.Error().Msgf("cant instantiate NewURLManagerServiceFromFullURL: %w", err)
		shutdown(&loggerInstance)
	}

	mux := http.GetRouter(&loggerInstance, shortcutService, urlManagerService)

	if err := nethttp.ListenAndServe(options.ServerListenAddr, mux); err != nil {
		loggerInstance.Error().Msgf("can't start server: ", err)
		shutdown(&loggerInstance)
	}
}

func shutdown(logger *zerolog.Logger) {
	logger.Info().Msg("shutting down")
	os.Exit(1)
}
