package http

import (
	"encoding/json"
	"fmt"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/config"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/http/api"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/logger"
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"io"
	"net/http"
)

func ShortenerApiHandler(res http.ResponseWriter, req *http.Request) {

	mimeType := ExtractMIMETypeFromStr(req.Header.Get(headers.ContentType))
	if mimeType != mimetype.ApplicationJSON {
		httpError(res, fmt.Sprintf("Content-type \"%s\" not allowed", mimeType))
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
		httpError(res, err.Error())
		return
	}

	shortcut, err := app.MakeShortcut(shortenerRequest.Url)
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
