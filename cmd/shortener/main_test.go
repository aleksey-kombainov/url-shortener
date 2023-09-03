package main

import (
	"github.com/aleksey-kombainov/url-shortener.git/pkg/config"
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

const (
	baseUri = `/`
)

var (
	testsShortener = []struct {
		//name string too lazy
		postData string
	}{
		{
			postData: `https://practicum.yandex.ru/`,
		},
		{
			postData: `http://ya.ru`,
		},
	}
)

func TestShortenerOK(t *testing.T) {
	shortcuts := make(map[string]string)

	for i, test := range testsShortener {
		t.Run(`Shortener test #`+strconv.Itoa(i), func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, baseUri, strings.NewReader(test.postData))
			request.Header.Add(headers.ContentType, mimetype.TextPlain)

			recorder := httptest.NewRecorder()
			getRouter().ServeHTTP(recorder, request)
			res := recorder.Result()

			assert.Equal(t, http.StatusCreated, res.StatusCode)
			assert.Equal(t, mimetype.TextPlain, res.Header.Get(headers.ContentType))

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.NotEmpty(t, string(resBody))

			shortcuts[strings.TrimPrefix(string(resBody), config.GetOptions().BaseUrl+baseUri)] = test.postData
		})
	}

	var i int = 0
	for shortcut, domain := range shortcuts {
		t.Run(`Expander test #`+strconv.Itoa(i), func(t *testing.T) {

			request := httptest.NewRequest(http.MethodGet, baseUri+shortcut, nil)

			recorder := httptest.NewRecorder()
			getRouter().ServeHTTP(recorder, request)
			res := recorder.Result()

			assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Empty(t, string(resBody))
			assert.Equal(t, domain, res.Header.Get(headers.Location))
		})
		i++
	}
}

func TestShortenerFailure(t *testing.T) {
	tests := []struct {
		// name string too lazy
		uri         string
		contentType string
		postData    string
	}{
		{
			uri:         "/xxx",
			contentType: mimetype.TextPlain,
			postData:    "practicum.yandex.ru",
		},
		{
			uri:         "/",
			contentType: mimetype.TextHTML,
			postData:    "practicum.yandex.ru",
		},
		{
			uri:         "/",
			contentType: mimetype.TextPlain,
			postData:    "",
		},
	}
	for i, test := range tests {
		t.Run(`Shortener test #`+strconv.Itoa(i), func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, test.uri, strings.NewReader(test.postData))
			request.Header.Add(headers.ContentType, test.contentType)

			recorder := httptest.NewRecorder()
			getRouter().ServeHTTP(recorder, request)
			res := recorder.Result()

			assert.Equal(t, http.StatusBadRequest, res.StatusCode)

			defer res.Body.Close()

			assert.NotEqual(t, mimetype.TextPlain, res.Header.Get(headers.ContentType))
		})
	}
}
