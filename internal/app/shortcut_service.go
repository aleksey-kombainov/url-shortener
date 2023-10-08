package app

import (
	"errors"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/random"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage"
	"github.com/rs/zerolog"
)

const (
	generatorIterationLimit = 1000
	shortcutLength          = 8
)

type ShortcutService struct {
	logger  *zerolog.Logger
	Storage *storage.ShortcutStorager
}

func NewShortcutService(logger *zerolog.Logger, storage *storage.ShortcutStorager) *ShortcutService {
	return &ShortcutService{
		logger:  logger,
		Storage: storage,
	}
}

func (s ShortcutService) MakeShortcut(url string) (shortcut string, err error) {
	if url == "" {
		return "", errors.New("url is empty")
	}
	isGenerated := false
	for i := 0; i < generatorIterationLimit; i++ {
		shortcut = random.GenString(shortcutLength)
		_, err := (*s.Storage).GetOriginalURLByShortcut(shortcut)
		if err != nil && errors.Is(err, storage.EntityNotFoundErr) { // err != nil - шорткат не найден
			isGenerated = true
			break
		} else if err != nil {
			s.logger.Error().Msgf("creating shortcut - error while getting original url: %w", err)
			return "", err
		}
	}
	if !isGenerated {
		return "", errors.New("generator limit exceeded")
	}
	err = (*s.Storage).CreateRecord(url, shortcut)
	return
}

func (s ShortcutService) GetOriginalURLByShortcut(shortcut string) (origURL string, err error) {
	return (*s.Storage).GetOriginalURLByShortcut(shortcut)
}

func (s ShortcutService) GetShortcutByOrenviginalURL(origURL string) (shortURL string, err error) {
	return (*s.Storage).GetShortcutByOriginalURL(origURL)
}
