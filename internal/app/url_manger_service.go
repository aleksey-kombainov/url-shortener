package app

import (
	"net/url"
	"strings"
)

type URLManagerService struct {
	Scheme  string
	BaseURL string
	BaseURI string
}

func NewURLManagerServiceFromFullURL(fullURL string) (*URLManagerService, error) {
	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}
	path := u.Path
	if path == "" {
		path = "/"
	}
	return &URLManagerService{
		Scheme:  u.Scheme,
		BaseURL: u.Host,
		BaseURI: path,
	}, nil
}

func (receiver URLManagerService) BuildFullURLByShortcut(shortcut string) string {
	return receiver.getBaseURL() + shortcut
}

func (receiver URLManagerService) GetShortcutFromFullURL(url string) string {
	return strings.TrimPrefix(url, receiver.getBaseURL())
}

func (receiver URLManagerService) GetShortcutFromURI(url string) string {
	return strings.TrimPrefix(url, receiver.BaseURI)
}

func (receiver URLManagerService) getBaseURL() string {
	return receiver.Scheme + "://" + receiver.BaseURL + receiver.BaseURI
}
