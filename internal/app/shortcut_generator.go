package app

import (
	"errors"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/random"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage"
)

const (
	generatorIterationLimit = 1000
	shortcutLength          = 8
)

func GenerateAndSaveRandomShortcut(url string) (string, error) {
	var shortcut string
	isGenerated := false
	for i := 0; i < generatorIterationLimit; i++ {
		shortcut = random.GenString(shortcutLength)
		if _, err := storage.ShortcutStorage.GetOriginalURLByShortcut(shortcut); err != nil {
			isGenerated = true
			break
		}
	}
	if !isGenerated {
		return "", errors.New("generator limit exceeded")
	}

	storage.ShortcutStorage.CreateRecord(url, shortcut)

	return shortcut, nil
}
