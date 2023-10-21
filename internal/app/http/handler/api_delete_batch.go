package handler

import (
	"encoding/json"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/rs/zerolog"
	"io"
	nethttp "net/http"
)

type DeleteBatchAPIHandler struct {
	logger          *zerolog.Logger
	shortcutService *app.ShortcutService
	urlService      *app.URLManagerService
}

func NewDeleteBatchAPIHandler(logger *zerolog.Logger, shortcutService *app.ShortcutService, urlService *app.URLManagerService) *DeleteBatchAPIHandler {
	return &DeleteBatchAPIHandler{logger: logger, shortcutService: shortcutService, urlService: urlService}
}

func (h DeleteBatchAPIHandler) ServeHTTP(res nethttp.ResponseWriter, req *nethttp.Request) {

	defer func() {
		if err := req.Body.Close(); err != nil {
			h.logger.Error().
				Msg("Can not close request.Body(): " + err.Error())
		}
	}()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.httpError(res, "can not read req body: "+err.Error(), nethttp.StatusInternalServerError)
		return
	}

	var request []string
	err = json.Unmarshal(body, &request)
	if err != nil {
		h.httpError(res, "Unmarshalling error: "+err.Error()+"; Body: "+string(body), nethttp.StatusInternalServerError)
		return
	}
	userID := getUserIDFromCtx(req.Context())
	h.logger.Debug().Msgf("delete batch len: %d; userID: %s", len(request), userID)

	go h.shortcutService.DeleteByIDsAndUser(request, userID)

	res.WriteHeader(nethttp.StatusAccepted)
}

func (h DeleteBatchAPIHandler) httpError(res nethttp.ResponseWriter, errStr string, httpStatus int) {
	h.logger.Error().
		Msg("http error: " + errStr)
	nethttp.Error(res, errStr, httpStatus)
}
