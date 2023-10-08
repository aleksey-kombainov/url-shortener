package memstorage

import (
	"context"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/entities"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage/storageerr"
)

type Storage struct {
	shortcutList []entities.Shortcut
	maxID        uint64
}

func New() *Storage {
	return &Storage{
		shortcutList: make([]entities.Shortcut, 0),
		maxID:        0,
	}
}

func (s *Storage) CreateRecord(origURL string, shortURL string) (err error) {
	s.maxID++
	rec := entities.Shortcut{
		ID:          s.maxID,
		ShortURL:    shortURL,
		OriginalURL: origURL,
	}
	if _, err = s.GetOriginalURLByShortcut(shortURL); err == nil {
		return storageerr.ErrNotUniqueShortcut
	} else if _, err = s.GetShortcutByOriginalURL(origURL); err == nil {
		return storageerr.ErrNotUniqueOriginalURL
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
	return "", storageerr.ErrEntityNotFound
}

func (s Storage) GetShortcutByOriginalURL(origURL string) (shortURL string, err error) {
	for _, sh := range s.shortcutList {
		if sh.OriginalURL == origURL {
			return sh.ShortURL, nil
		}
	}
	return "", storageerr.ErrEntityNotFound
}

func (*Storage) Close() (err error) {
	return nil
}

func (s *Storage) Ping(ctx context.Context) (err error) {
	return nil
}
