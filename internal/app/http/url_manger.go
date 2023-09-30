package http

import (
	"net/url"
	"strings"
)

type urlManager struct {
	Scheme  string
	BaseURL string
	BaseURI string
}

func NewURLManagerFromFullURL(fullURL string) *urlManager {
	u, err := url.Parse(fullURL)
	if err != nil {
		panic(err)
	}
	path := u.Path
	if path == "" {
		path = "/"
	}
	return &urlManager{
		Scheme:  u.Scheme,
		BaseURL: u.Host,
		BaseURI: path,
	}
}

func (receiver urlManager) BuildFullURLByShortcut(shortcut string) string {
	return receiver.getBaseURL() + shortcut
}

func (receiver urlManager) GetShortcutFromFullURL(url string) string {
	return strings.TrimPrefix(url, receiver.getBaseURL())
}

func (receiver urlManager) GetShortcutFromURI(url string) string {
	return strings.TrimPrefix(url, receiver.BaseURI)
}

func (receiver urlManager) getBaseURL() string {
	return receiver.Scheme + "://" + receiver.BaseURL + receiver.BaseURI
}
