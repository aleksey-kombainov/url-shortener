package http

import (
	"encoding/json"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/config"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/http/api"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/logger"
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"io"
	"net/http"
)

func ShortenerAPIHandler(res http.ResponseWriter, req *http.Request) {

	if !IsHeaderContainsMIMEType(req.Header.Values(headers.ContentType), mimetype.ApplicationJSON) {
		httpError(res, "Content-type not allowed")
		return
	}
	defer func() {
		if err := req.Body.Close(); err != nil {
			logger.Logger.Error().
				Msg("Can not close request.Body(): " + err.Error())
		}
	}()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		httpError(res, err.Error())
		return
	}

	shortenerRequest := api.ShortenerRequest{}
	err = json.Unmarshal(body, &shortenerRequest)
	if err != nil {
		httpError(res, "Unmarshalling error: "+err.Error()+"; Body: "+string(body))
		return
	}

	shortcut, err := app.MakeShortcut(shortenerRequest.URL)
	if err != nil {
		httpError(res, err.Error())
		return
	}

	url := NewURLManagerFromFullURL(config.GetOptions().BaseURL).BuildFullURLByShortcut(shortcut)
	response, err := json.Marshal(api.ShortenerResponse{Result: url})
	if err != nil {
		httpError(res, err.Error())
	}

	res.Header().Add(headers.ContentType, mimetype.ApplicationJSON)
	res.WriteHeader(http.StatusCreated)

	if _, err := res.Write(response); err != nil {
		httpError(res, err.Error())
	}
}
