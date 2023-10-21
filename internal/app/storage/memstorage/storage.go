package memstorage

import (
	"context"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/entities"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/interfaces"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/model"
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

func (s *Storage) CreateRecord(ctx context.Context, origURL string, shortURL string, userID string) (shortcut entities.Shortcut, err error) {
	s.maxID++
	shortcut = entities.Shortcut{
		ID:          s.maxID,
		ShortURL:    shortURL,
		OriginalURL: origURL,
		UserID:      userID,
	}
	if _, err = s.GetOriginalURLByShortcut(ctx, shortURL); err == nil {
		return shortcut, storageerr.ErrNotUniqueShortcut
	} else if _, err = s.GetShortcutByOriginalURL(ctx, origURL); err == nil {
		return shortcut, storageerr.ErrNotUniqueOriginalURL
	}
	s.shortcutList = append(s.shortcutList, shortcut)
	return
}

func (s Storage) GetOriginalURLByShortcut(ctx context.Context, shortURL string) (shortcut entities.Shortcut, err error) {
	for _, shortcut = range s.shortcutList {
		if shortcut.ShortURL == shortURL {
			return shortcut, nil
		}
	}
	return shortcut, storageerr.ErrEntityNotFound
}

func (s Storage) GetShortcutByOriginalURL(ctx context.Context, origURL string) (shortcut entities.Shortcut, err error) {
	for _, shortcut = range s.shortcutList {
		if shortcut.OriginalURL == origURL {
			return shortcut, nil
		}
	}
	return shortcut, storageerr.ErrEntityNotFound
}

func (*Storage) Close(ctx context.Context) (err error) {
	return nil
}

func (s *Storage) Ping(ctx context.Context) (err error) {
	return nil
}

func (s *Storage) NewBatch(ctx context.Context) (interfaces.ShortcutStorager, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) CreateRecordBatch(ctx context.Context, origURL string, shortURL string, userID string) (shortcut entities.Shortcut, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) CommitBatch(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) RollbackBatch(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetShortcutsByUser(ctx context.Context, userID string) (shortcuts []entities.Shortcut, err error) {
	for _, sh := range s.shortcutList {
		if sh.UserID == userID {
			shortcuts = append(shortcuts, sh)
		}
	}
	return
}

func (s Storage) DeleteByShortcutsForUser(ctx context.Context, shortcuts []string, userID string) (err error) {
	return
}

func (s Storage) DeleteByShortcutsAndUser(ctx context.Context, deleteTasks []model.DeleteTask) (err error) {
	panic("implement me")
}
