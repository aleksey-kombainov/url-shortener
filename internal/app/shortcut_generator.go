package app

import (
	"errors"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/memstorage"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/random"
)

const (
	generatorIterationLimit = 1000
	shortcutLength          = 8
)

func GetAndSaveUniqueShortcut(url string, storage memstorage.Storager) (string, error) {
	var shortcut string
	isGenerated := false
	for i := 0; i < generatorIterationLimit; i++ {
		shortcut = random.GenString(shortcutLength)
		if _, err := storage.GetValueByKey(shortcut); err != nil {
			isGenerated = true
			break
		}
	}
	if !isGenerated {
		return "", errors.New("generator limit exceeded")
	}
	storage.Put(shortcut, url)
	return shortcut, nil
}
