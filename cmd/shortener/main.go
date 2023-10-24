package main

import (
	"context"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/config"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/http"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/interfaces"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/logger"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage"
	"github.com/rs/zerolog"
	nethttp "net/http"
	"os"
	"time"
)

func main() {
	var err error

	loggerInstance := logger.GetLogger(zerolog.DebugLevel)

	options := config.GetOptions(loggerInstance)
	loggerInstance.Debug().Msgf("options: %+v", options)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var storageInstance interfaces.ShortcutStorager
	if options.DatabaseDsn != "" {
		storageInstance, err = storage.ShortcutStorageFactory(ctx, &loggerInstance, storage.TypeDB, options.DatabaseDsn)
	} else if options.FileStoragePath != "" {
		storageInstance, err = storage.ShortcutStorageFactory(ctx, &loggerInstance, storage.TypeFile, options.FileStoragePath)
	} else {
		storageInstance, err = storage.ShortcutStorageFactory(ctx, &loggerInstance, storage.TypeMemory, "")
	}
	if err != nil {
		loggerInstance.Error().Msgf("cant get storage: %s", err)
		shutdown(&loggerInstance)
	} else {
		defer func() {
			if err := storageInstance.Close(context.TODO()); err != nil {
				loggerInstance.Error().Msgf("can't close storage: %s", err)
			}
		}()
	}

	loggerInstance.Info().Msg("Starting server")

	shortcutService := app.NewShortcutService(&loggerInstance, storageInstance)
	urlManagerService, err := app.NewURLManagerServiceFromFullURL(options.BaseURL)
	if err != nil {
		loggerInstance.Error().Msgf("can' t instantiate NewURLManagerServiceFromFullURL: %s", err)
		shutdown(&loggerInstance)
	}

	mux := http.GetRouter(&loggerInstance, shortcutService, urlManagerService)

	if err := nethttp.ListenAndServe(options.ServerListenAddr, mux); err != nil {
		loggerInstance.Error().Msgf("can't start server: %s", err)
		shutdown(&loggerInstance)
	}
}

func shutdown(logger *zerolog.Logger) {
	logger.Info().Msg("shutting down")
	os.Exit(1)
}
