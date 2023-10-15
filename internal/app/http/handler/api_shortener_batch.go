package handler

import (
	"context"
	"encoding/json"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/model"
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"github.com/rs/zerolog"
	"io"
	"net/http"
)

type ShortenerBatchAPIHandler struct {
	logger          *zerolog.Logger
	shortcutService *app.ShortcutService
	urlService      *app.URLManagerService
}

func NewShortenerBatchAPIHandler(logger *zerolog.Logger, shortcutService *app.ShortcutService, urlService *app.URLManagerService) *ShortenerBatchAPIHandler {
	return &ShortenerBatchAPIHandler{logger: logger, shortcutService: shortcutService, urlService: urlService}
}

func (h ShortenerBatchAPIHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	defer func() {
		if err := req.Body.Close(); err != nil {
			h.logger.Error().
				Msg("Can not close request.Body(): " + err.Error())
		}
	}()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.httpError(res, "can not read req body: "+err.Error(), http.StatusBadRequest)
		return
	}

	var shortenerRequest []model.ShortenerBatchRecordRequest
	err = json.Unmarshal(body, &shortenerRequest)
	if err != nil {
		h.httpError(res, "Unmarshalling error: "+err.Error()+"; Body: "+string(body), http.StatusBadRequest)
		return
	}

	responseJSON, err := h.shortcutService.MakeShortcutBatch(context.TODO(), shortenerRequest)
	if err != nil {
		h.httpError(res, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(responseJSON)
	if err != nil {
		h.httpError(res, "Marshalling error: "+err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Add(headers.ContentType, mimetype.ApplicationJSON)
	res.WriteHeader(http.StatusCreated)

	if _, err := res.Write(response); err != nil {
		h.httpError(res, "Writing response error: "+err.Error(), http.StatusBadRequest)
		return
	}
}

func (h ShortenerBatchAPIHandler) httpError(res http.ResponseWriter, errStr string, httpStatus int) {
	h.logger.Error().
		Msg("http error: " + errStr)
	http.Error(res, errStr, httpStatus)
}
