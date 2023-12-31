package main

import (
	"testing"
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
	shortcuts = make(map[string]string)
)

func TestShortenerOK(t *testing.T) {
	//initGlobals()
	//options = config.GetOptions()
	//for i, test := range testsShortener {
	//	t.Run(`Shortener test #`+strconv.Itoa(i), func(t *testing.T) {
	//		request := httptest.NewRequest(nethttp.MethodPost, http.NewURLManagerFromFullURL(config.GetOptions().BaseURL).BaseURI, strings.NewReader(test.postData))
	//		request.Header.Add(headers.ContentType, mimetype.TextPlain)
	//
	//		recorder := httptest.NewRecorder()
	//		http.GetRouter().ServeHTTP(recorder, request)
	//		res := recorder.Result()
	//		defer res.Body.Close()
	//
	//		assert.Equal(t, nethttp.StatusCreated, res.StatusCode)
	//		assert.True(t, http.IsHeaderContainsMIMETypes(res.Header.Values(headers.ContentType), []string{mimetype.TextPlain}))
	//
	//		resBody, err := io.ReadAll(res.Body)
	//		require.NoError(t, err)
	//		assert.NotEmpty(t, string(resBody))
	//
	//		shortcuts[http.NewURLManagerFromFullURL(config.GetOptions().BaseURL).GetShortcutFromFullURL(string(resBody))] = test.postData
	//	})
	//}
	//TestShortenerAPIOK(t)
}

func TestShortenerAPIOK(t *testing.T) {
	//options = config.GetOptions()
	//reqURL := http.NewURLManagerFromFullURL(config.GetOptions().BaseURL).BaseURI + "model/shorten"
	//for i, test := range testsShortener {
	//	t.Run(`Shortener model test #`+strconv.Itoa(i), func(t *testing.T) {
	//		reqStr, _ := json.Marshal(model.ShortenerRequest{URL: test.postData})
	//		request := httptest.NewRequest(nethttp.MethodPost, reqURL, bytes.NewReader(reqStr))
	//		request.Header.Add(headers.ContentType, mimetype.ApplicationJSON)
	//
	//		recorder := httptest.NewRecorder()
	//		http.GetRouter().ServeHTTP(recorder, request)
	//		res := recorder.Result()
	//		defer res.Body.Close()
	//
	//		assert.Equal(t, nethttp.StatusCreated, res.StatusCode)
	//		assert.True(t, http.IsHeaderContainsMIMETypes(res.Header.Values(headers.ContentType), []string{mimetype.ApplicationJSON}))
	//
	//		resBody, err := io.ReadAll(res.Body)
	//		require.NoError(t, err)
	//		assert.NotEmpty(t, string(resBody))
	//
	//		sresp := &model.ShortenerResponse{}
	//		json.Unmarshal(resBody, sresp)
	//		sc := http.NewURLManagerFromFullURL(config.GetOptions().BaseURL).GetShortcutFromFullURL(sresp.Result)
	//		shortcuts[sc] = test.postData
	//	})
	//}
	//expanderOK(t)
}

func expanderOK(t *testing.T) {
	//var i = 0
	//for shortcut, domain := range shortcuts {
	//	t.Run(`Expander test #`+strconv.Itoa(i), func(t *testing.T) {
	//
	//		request := httptest.NewRequest(nethttp.MethodGet, http.NewURLManagerFromFullURL(config.GetOptions().BaseURL).BaseURI+shortcut, nil)
	//
	//		recorder := httptest.NewRecorder()
	//		http.GetRouter().ServeHTTP(recorder, request)
	//		res := recorder.Result()
	//		defer res.Body.Close()
	//
	//		assert.Equal(t, nethttp.StatusTemporaryRedirect, res.StatusCode)
	//
	//		resBody, err := io.ReadAll(res.Body)
	//
	//		require.NoError(t, err)
	//		assert.Empty(t, string(resBody))
	//		assert.Equal(t, domain, res.Header.Get(headers.Location))
	//	})
	//	i++
	//}
}

func TestShortenerFailure(t *testing.T) {
	//tests := []struct {
	//	// name string too lazy
	//	uri         string
	//	contentType string
	//	postData    string
	//}{
	//	{
	//		uri:         "/xxx",
	//		contentType: mimetype.TextPlain,
	//		postData:    "practicum.yandex.ru",
	//	},
	//	{
	//		uri:         "/",
	//		contentType: mimetype.TextHTML,
	//		postData:    "practicum.yandex.ru",
	//	},
	//	{
	//		uri:         "/",
	//		contentType: mimetype.TextPlain,
	//		postData:    "",
	//	},
	//}
	//for i, test := range tests {
	//	t.Run(`Shortener test #`+strconv.Itoa(i), func(t *testing.T) {
	//
	//		request := httptest.NewRequest(nethttp.MethodPost, test.uri, strings.NewReader(test.postData))
	//		request.Header.Add(headers.ContentType, test.contentType)
	//
	//		recorder := httptest.NewRecorder()
	//		http.GetRouter().ServeHTTP(recorder, request)
	//		res := recorder.Result()
	//		defer res.Body.Close()
	//
	//		assert.Equal(t, nethttp.StatusBadRequest, res.StatusCode)
	//
	//		assert.NotEqual(t, mimetype.TextPlain, res.Header.Get(headers.ContentType))
	//	})
	//}
}
