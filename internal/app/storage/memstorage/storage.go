package memstorage

import (
	"context"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/entities"
)

type Storage struct {
	shortcutList      []entities.Shortcut
	maxID             uint64
	entityNotFoundErr error
}

func New(entityNotFoundErr error) *Storage {
	return &Storage{
		shortcutList:      make([]entities.Shortcut, 0),
		maxID:             0,
		entityNotFoundErr: entityNotFoundErr,
	}
}

func (s *Storage) CreateRecord(origURL string, shortURL string) (err error) {
	s.maxID++
	rec := entities.Shortcut{
		ID:          s.maxID,
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
	return "", s.entityNotFoundErr
}

func (s Storage) GetShortcutByOriginalURL(origURL string) (shortURL string, err error) {
	for _, sh := range s.shortcutList {
		if sh.OriginalURL == origURL {
			return sh.ShortURL, nil
		}
	}
	return "", s.entityNotFoundErr
}

func (*Storage) Close() (err error) {
	return nil
}

func (s *Storage) Ping(ctx context.Context) (err error) {
	return nil
}
