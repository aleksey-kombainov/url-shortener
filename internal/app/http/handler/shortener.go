package handler

import (
	"errors"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage/storageerr"
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"strings"
)

type ShortenerHandler struct {
	logger          *zerolog.Logger
	shortcutService *app.ShortcutService
	urlService      *app.URLManagerService
}

func NewShortenerHandler(logger *zerolog.Logger, shortcutService *app.ShortcutService, urlService *app.URLManagerService) *ShortenerHandler {
	return &ShortenerHandler{logger: logger, shortcutService: shortcutService, urlService: urlService}
}

func (h ShortenerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			h.logger.Error().
				Msg("Can not close request.Body(): " + err.Error())
		}
	}()
	url, err := io.ReadAll(req.Body)
	if err != nil {
		h.httpError(res, err.Error())
		return
	}

	httpStatus := http.StatusCreated
	shortcut, err := h.shortcutService.MakeShortcut(strings.TrimSpace(string(url)), getUserIDFromCtx(req.Context()))
	if errors.Is(err, storageerr.ErrNotUniqueOriginalURL) {
		httpStatus = http.StatusConflict
	} else if err != nil {
		h.httpError(res, "MakeShortcut: "+err.Error())
		return
	}
	res.Header().Add(headers.ContentType, mimetype.TextPlain)
	res.WriteHeader(httpStatus)
	if _, err := res.Write([]byte(h.urlService.BuildFullURLByShortcut(shortcut))); err != nil {
		h.logger.Error().
			Msg("Can not Write response: " + err.Error())
	}
}

func (h ShortenerHandler) httpError(res http.ResponseWriter, errStr string) {
	h.logger.Error().
		Msg("http error: " + errStr)
	http.Error(res, errStr, http.StatusBadRequest)
}
