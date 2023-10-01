package http

import (
	"net/url"
	"strings"
)

type URLManager struct {
	Scheme  string
	BaseURL string
	BaseURI string
}

func NewURLManagerFromFullURL(fullURL string) *URLManager {
	u, err := url.Parse(fullURL)
	if err != nil {
		panic(err)
	}
	path := u.Path
	if path == "" {
		path = "/"
	}
	return &URLManager{
		Scheme:  u.Scheme,
		BaseURL: u.Host,
		BaseURI: path,
	}
}

func (receiver URLManager) BuildFullURLByShortcut(shortcut string) string {
	return receiver.getBaseURL() + shortcut
}

func (receiver URLManager) GetShortcutFromFullURL(url string) string {
	return strings.TrimPrefix(url, receiver.getBaseURL())
}

func (receiver URLManager) GetShortcutFromURI(url string) string {
	return strings.TrimPrefix(url, receiver.BaseURI)
}

func (receiver URLManager) getBaseURL() string {
	return receiver.Scheme + "://" + receiver.BaseURL + receiver.BaseURI
}
