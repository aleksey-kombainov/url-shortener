package http

import (
	"net/url"
	"strings"
)

type UrlManager struct {
	Scheme  string
	BaseURL string
	BaseURI string
}

func NewURLManagerFromFullURL(fullURL string) *UrlManager {
	u, err := url.Parse(fullURL)
	if err != nil {
		panic(err)
	}
	path := u.Path
	if path == "" {
		path = "/"
	}
	return &UrlManager{
		Scheme:  u.Scheme,
		BaseURL: u.Host,
		BaseURI: path,
	}
}

func (receiver UrlManager) BuildFullURLByShortcut(shortcut string) string {
	return receiver.getBaseURL() + shortcut
}

func (receiver UrlManager) GetShortcutFromFullURL(url string) string {
	return strings.TrimPrefix(url, receiver.getBaseURL())
}

func (receiver UrlManager) GetShortcutFromURI(url string) string {
	return strings.TrimPrefix(url, receiver.BaseURI)
}

func (receiver UrlManager) getBaseURL() string {
	return receiver.Scheme + "://" + receiver.BaseURL + receiver.BaseURI
}
