package memstorage

import (
	"errors"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/entities"
)

type Storage struct {
	shortcutList []entities.Shortcut
	maxId        uint64
}

func New() *Storage {
	return &Storage{
		shortcutList: make([]entities.Shortcut, 0),
		maxId:        0,
	}
}

func (s *Storage) CreateRecord(origURL string, shortURL string) (err error) {
	s.maxId++
	rec := entities.Shortcut{
		ID:          s.maxId,
		ShortURL:    shortURL,
		OriginalURL: origURL,
	}
	s.shortcutList = append(s.shortcutList, rec)
	return nil
}

func (s Storage) GetOriginalURLByShortcut(shortURL string) (origURL string, err error) {
	for _, sh := range s.shortcutList {
		if sh.ShortURL == shortURL {
			return sh.OriginalURL, nil
		}
	}
	return "", errors.New("shortcut not found")
}

func (s Storage) GetShortcutByOriginalURL(origURL string) (shortURL string, err error) {
	for _, sh := range s.shortcutList {
		if sh.OriginalURL == origURL {
			return sh.ShortURL, nil
		}
	}
	return "", errors.New("original url not found")
}

func (*Storage) Close() (err error) {
	return nil
}
