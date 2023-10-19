package handler

import (
	"encoding/json"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"github.com/rs/zerolog"
	"net/http"
)

type UserURLsAPIHandler struct {
	logger          *zerolog.Logger
	shortcutService *app.ShortcutService
	urlService      *app.URLManagerService
}

func NewUserURLsAPIHandler(logger *zerolog.Logger, shortcutService *app.ShortcutService, urlService *app.URLManagerService) *UserURLsAPIHandler {
	return &UserURLsAPIHandler{logger: logger, shortcutService: shortcutService, urlService: urlService}
}

func (h UserURLsAPIHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	userID := getUserIDFromCtx(req.Context())
	if userID == "" {
		h.httpError(res, "", http.StatusUnauthorized)
		return
	}

	shortcuts, err := h.shortcutService.GetShortcutsByUser(userID)
	if err != nil {
		h.httpError(res, "", http.StatusInternalServerError)
		return
	}
	if len(shortcuts) == 0 {
		h.httpError(res, "", http.StatusNoContent)
		return
	}
	response, err := json.Marshal(shortcuts)
	if err != nil {
		h.httpError(res, "Marshalling error: "+err.Error(), http.StatusInternalServerError)
	}

	res.Header().Add(headers.ContentType, mimetype.ApplicationJSON)
	res.WriteHeader(http.StatusOK)

	if _, err := res.Write(response); err != nil {
		h.httpError(res, "Writing response error: "+err.Error(), http.StatusBadRequest)
	}
}

func (h UserURLsAPIHandler) httpError(res http.ResponseWriter, errStr string, httpStatus int) {
	h.logger.Error().
		Msg("http error: " + errStr)
	http.Error(res, errStr, httpStatus)
}
