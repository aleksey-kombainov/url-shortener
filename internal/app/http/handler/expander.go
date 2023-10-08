package handler

import (
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/go-http-utils/headers"
	"github.com/rs/zerolog"
	"net/http"
)

type ExpanderHandler struct {
	logger          *zerolog.Logger
	shortcutService *app.ShortcutService
	urlService      *app.URLManagerService
}

func NewExpanderHandler(logger *zerolog.Logger, shortcutService *app.ShortcutService, urlService *app.URLManagerService) *ExpanderHandler {
	return &ExpanderHandler{logger: logger, shortcutService: shortcutService, urlService: urlService}
}

func (h ExpanderHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	shortcut := h.urlService.GetShortcutFromURI(req.RequestURI)
	if len(shortcut) == 0 {
		h.httpError(res, "invalid shortcut")
		return
	}
	// @todo errors
	url, err := h.shortcutService.GetOriginalURLByShortcut(shortcut)
	if err != nil {
		h.httpError(res, "shortcut not found")
		return
	}
	res.Header().Add(headers.Location, url) // @todo проверить редирект на самого себя
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h ExpanderHandler) httpError(res http.ResponseWriter, errStr string) {
	h.logger.Error().
		Msg("http error: " + errStr)
	http.Error(res, errStr, http.StatusBadRequest)
}
