package handler

import (
	"encoding/json"
	"errors"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/http/api"
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

func (s ShortenerAPIHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	defer func() {
		if err := req.Body.Close(); err != nil {
			s.logger.Error().
				Msg("Can not close request.Body(): " + err.Error())
		}
	}()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		s.httpError(res, "can not read req body: "+err.Error())
		return
	}

	shortenerRequest := api.ShortenerRequest{}
	err = json.Unmarshal(body, &shortenerRequest)
	if err != nil {
		s.httpError(res, "Unmarshalling error: "+err.Error()+"; Body: "+string(body))
		return
	}

	httpStatus := http.StatusCreated
	shortcut, err := s.shortcutService.MakeShortcut(strings.TrimSpace(shortenerRequest.URL))
	if errors.Is(err, storageerr.ErrNotUniqueOriginalURL) {
		httpStatus = http.StatusConflict
	} else if err != nil {
		s.httpError(res, "MakeShortcut: "+err.Error())
		return
	}

	url := s.urlService.BuildFullURLByShortcut(shortcut)
	response, err := json.Marshal(api.ShortenerResponse{Result: url})
	if err != nil {
		s.httpError(res, "Marshalling error: "+err.Error())
	}

	res.Header().Add(headers.ContentType, mimetype.ApplicationJSON)
	res.WriteHeader(httpStatus)

	if _, err := res.Write(response); err != nil {
		s.httpError(res, "Writing response error: "+err.Error())
	}
}

func (s ShortenerAPIHandler) httpError(res http.ResponseWriter, errStr string) {
	s.logger.Error().
		Msg("http error: " + errStr)
	http.Error(res, errStr, http.StatusBadRequest)
}
