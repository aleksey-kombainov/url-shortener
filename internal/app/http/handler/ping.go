package handler

import (
	"context"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/interfaces"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type PingHandler struct {
	logger  *zerolog.Logger
	storage *interfaces.ShortcutStorager
}

func NewPingHandler(logger *zerolog.Logger, storage interfaces.ShortcutStorager) *PingHandler {
	return &PingHandler{logger: logger, storage: &storage}
}

func (h PingHandler) ServeHTTP(res http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := (*h.storage).Ping(ctx); err != nil {
		(*h.logger).Error().Msg("unable to ping: " + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
