package handler

import (
	"context"
	"encoding/json"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/model"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/user"
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"github.com/rs/zerolog"
	"io"
	nethttp "net/http"
)

type ShortenerBatchAPIHandler struct {
	logger          *zerolog.Logger
	shortcutService *app.ShortcutService
	urlService      *app.URLManagerService
}

func NewShortenerBatchAPIHandler(logger *zerolog.Logger, shortcutService *app.ShortcutService, urlService *app.URLManagerService) *ShortenerBatchAPIHandler {
	return &ShortenerBatchAPIHandler{logger: logger, shortcutService: shortcutService, urlService: urlService}
}

func (h ShortenerBatchAPIHandler) ServeHTTP(res nethttp.ResponseWriter, req *nethttp.Request) {

	defer func() {
		if err := req.Body.Close(); err != nil {
			h.logger.Error().
				Msg("Can not close request.Body(): " + err.Error())
		}
	}()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.httpError(res, "can not read req body: "+err.Error(), nethttp.StatusBadRequest)
		return
	}

	var shortenerRequest []model.ShortenerBatchRecordRequest
	err = json.Unmarshal(body, &shortenerRequest)
	if err != nil {
		h.httpError(res, "Unmarshalling error: "+err.Error()+"; Body: "+string(body), nethttp.StatusBadRequest)
		return
	}

	userID := req.Context().Value(user.CtxUserIDKey).(string)
	shortenerBatchRecordResponses, err := h.shortcutService.MakeShortcutBatch(context.TODO(), shortenerRequest, userID)
	if err != nil {
		h.httpError(res, err.Error(), nethttp.StatusBadRequest)
		return
	}

	for idx, respRecord := range shortenerBatchRecordResponses {
		shortenerBatchRecordResponses[idx].ShortURL = h.urlService.BuildFullURLByShortcut(respRecord.ShortURL)
	}

	response, err := json.Marshal(shortenerBatchRecordResponses)
	if err != nil {
		h.httpError(res, "Marshalling error: "+err.Error(), nethttp.StatusBadRequest)
		return
	}

	res.Header().Add(headers.ContentType, mimetype.ApplicationJSON)
	res.WriteHeader(nethttp.StatusCreated)

	if _, err := res.Write(response); err != nil {
		h.httpError(res, "Writing response error: "+err.Error(), nethttp.StatusBadRequest)
		return
	}
}

func (h ShortenerBatchAPIHandler) httpError(res nethttp.ResponseWriter, errStr string, httpStatus int) {
	h.logger.Error().
		Msg("http error: " + errStr)
	nethttp.Error(res, errStr, httpStatus)
}
