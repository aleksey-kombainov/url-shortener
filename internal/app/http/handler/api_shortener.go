package handler

import (
	"encoding/json"
	"errors"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/model"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage/storageerr"
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"strings"
)

type ShortenerAPIHandler struct {
	logger          *zerolog.Logger
	shortcutService *app.ShortcutService
	urlService      *app.URLManagerService
}

func NewShortenerAPIHandler(logger *zerolog.Logger, shortcutService *app.ShortcutService, urlService *app.URLManagerService) *ShortenerAPIHandler {
	return &ShortenerAPIHandler{logger: logger, shortcutService: shortcutService, urlService: urlService}
}

func (h ShortenerAPIHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	defer func() {
		if err := req.Body.Close(); err != nil {
			h.logger.Error().
				Msg("Can not close request.Body(): " + err.Error())
		}
	}()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.httpError(res, "can not read req body: "+err.Error())
		return
	}

	shortenerRequest := model.ShortenerRequest{}
	err = json.Unmarshal(body, &shortenerRequest)
	if err != nil {
		h.httpError(res, "Unmarshalling error: "+err.Error()+"; Body: "+string(body))
		return
	}

	httpStatus := http.StatusCreated
	shortcut, err := h.shortcutService.MakeShortcut(strings.TrimSpace(shortenerRequest.URL), getUserIDFromCtx(req.Context()))
	if errors.Is(err, storageerr.ErrNotUniqueOriginalURL) {
		httpStatus = http.StatusConflict
	} else if err != nil {
		h.httpError(res, "MakeShortcut: "+err.Error())
		return
	}

	url := h.urlService.BuildFullURLByShortcut(shortcut.ShortURL)
	response, err := json.Marshal(model.ShortenerResponse{Result: url})
	if err != nil {
		h.httpError(res, "Marshalling error: "+err.Error())
	}

	res.Header().Add(headers.ContentType, mimetype.ApplicationJSON)
	res.WriteHeader(httpStatus)

	if _, err := res.Write(response); err != nil {
		h.httpError(res, "Writing response error: "+err.Error())
	}
}

func (h ShortenerAPIHandler) httpError(res http.ResponseWriter, errStr string) {
	h.logger.Error().
		Msg("http error: " + errStr)
	http.Error(res, errStr, http.StatusBadRequest)
}
