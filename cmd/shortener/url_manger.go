package main

import (
	"net/url"
	"strings"
)

type urlManager struct {
	scheme  string
	baseURL string
	baseURI string
}

//func newURLManager(scheme string, baseURL string, baseURI string) *urlManager {
//	return &urlManager{
//		scheme:  scheme,
//		baseURL: baseURL,
//		baseURI: baseURI,
//	}
//}

func newURLManagerFromFullURL(fullURL string) *urlManager {
	u, err := url.Parse(fullURL)
	if err != nil {
		panic(err)
	}
	path := u.Path
	if path == "" {
		path = "/"
	}
	return &urlManager{
		scheme:  u.Scheme,
		baseURL: u.Host,
		baseURI: path,
	}
}

func (receiver urlManager) buildFullURLByShortcut(shortcut string) string {
	return receiver.getBaseURL() + shortcut
}

func (receiver urlManager) getShortcutFromFullURL(url string) string {
	return strings.TrimPrefix(url, receiver.getBaseURL())
}

func (receiver urlManager) getShortcutFromURI(url string) string {
	return strings.TrimPrefix(url, receiver.baseURI)
}

func (receiver urlManager) getBaseURL() string {
	return receiver.scheme + "://" + receiver.baseURL + receiver.baseURI
}
