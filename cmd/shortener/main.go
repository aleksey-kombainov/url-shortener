package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/config"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/http"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/interfaces"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/logger"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage"
	"github.com/rs/zerolog"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
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
	// https://github.com/go-chi/chi/blob/master/_examples/graceful/main.go
	// The HTTP Server
	server := &nethttp.Server{Addr: options.ServerListenAddr, Handler: mux}
	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				loggerInstance.Fatal().Msg("graceful shutdown timed out.. forcing exit.")
			}
		}()
		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			loggerInstance.Fatal().Msg(err.Error())
		}
		fmt.Println("hurra")
		serverStopCtx()
	}()

	// Run the server
	err = server.ListenAndServe()
	if err != nil && errors.Is(err, nethttp.ErrServerClosed) {
		loggerInstance.Error().Msg("1: " + err.Error())
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

func shutdown(logger *zerolog.Logger) {
	logger.Info().Msg("shutting down")
	os.Exit(1)
}
