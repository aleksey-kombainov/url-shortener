package app

import (
	"errors"
	"strings"
)

func MakeShortcut(url string) (shortcut string, err error) {
	urlStr := strings.TrimSpace(url) // @todo валидация url
	if urlStr == "" {
		return "", errors.New("url is empty")
	}

	shortcut, err = GenerateAndSaveRandomShortcut(url)
	if err != nil {
		return "", err
	}
	return
}
