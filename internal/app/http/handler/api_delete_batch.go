package handler

import (
	"context"
	"encoding/json"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/model"
	"github.com/rs/zerolog"
	"io"
	nethttp "net/http"
	"time"
)

type DeleteBatchAPIHandler struct {
	logger          *zerolog.Logger
	shortcutService *app.ShortcutService
	urlService      *app.URLManagerService
	deleteTasksChan chan model.DeleteTask
}

func NewDeleteBatchAPIHandler(logger *zerolog.Logger, shortcutService *app.ShortcutService, urlService *app.URLManagerService) *DeleteBatchAPIHandler {
	h := &DeleteBatchAPIHandler{
		logger:          logger,
		shortcutService: shortcutService,
		urlService:      urlService,
		deleteTasksChan: make(chan model.DeleteTask, 1024),
	}
	go h.flushDeleteQueue()
	return h
}

func (h DeleteBatchAPIHandler) ServeHTTP(res nethttp.ResponseWriter, req *nethttp.Request) {

	defer func() {
		if err := req.Body.Close(); err != nil {
			h.logger.Error().
				Msg("Can not close requestData.Body(): " + err.Error())
		}
	}()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.httpError(res, "can not read req body: "+err.Error(), nethttp.StatusInternalServerError)
		return
	}

	var requestData []string
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		h.httpError(res, "Unmarshalling error: "+err.Error()+"; Body: "+string(body), nethttp.StatusInternalServerError)
		return
	}
	userID := getUserIDFromCtx(req.Context())

	h.deleteTasksChan <- model.DeleteTask{
		UserID:    userID,
		ShortURLs: requestData,
	}

	res.WriteHeader(nethttp.StatusAccepted)
}

func (h DeleteBatchAPIHandler) httpError(res nethttp.ResponseWriter, errStr string, httpStatus int) {
	h.logger.Error().
		Msg("http error: " + errStr)
	nethttp.Error(res, errStr, httpStatus)
}

func (h DeleteBatchAPIHandler) flushDeleteQueue() {
	// будем удалять накопленные за последние 10 секунд
	ticker := time.NewTicker(10 * time.Second)

	var deleteTasks []model.DeleteTask

	for {
		select {
		case msg := <-h.deleteTasksChan:
			// добавим сообщение в слайс для последующего сохранения
			deleteTasks = append(deleteTasks, msg)
		case <-ticker.C:
			// подождём, пока придёт хотя бы одно сообщение
			if len(deleteTasks) == 0 {
				continue
			}
			h.logger.Debug().Msgf("delete batch len: %d", len(deleteTasks))
			err := h.shortcutService.DeleteByShortcutsAndUser(context.Background(), deleteTasks)
			if err != nil {
				h.logger.Error().Msgf("cannot save deleteTasks: %s", err.Error())
				// не будем стирать сообщения, попробуем отправить их чуть позже
				continue
			}
			// сотрём успешно отосланные сообщения
			deleteTasks = nil
		}
	}
}
