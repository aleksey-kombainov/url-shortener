package memstorage

import (
	"context"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/entities"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/interfaces"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/model"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage/storageerr"
	"sync"
	"sync/atomic"
)

type Storage struct {
	shortcutList        []entities.Shortcut
	maxID               atomic.Uint64
	shortcutListRWMutex sync.RWMutex
	batch               []entities.Shortcut
}

func New() *Storage {
	return &Storage{
		shortcutList: make([]entities.Shortcut, 0),
	}
}

func (s *Storage) CreateRecord(ctx context.Context, origURL string, shortURL string, userID string) (shortcut entities.Shortcut, err error) {
	nextID := s.incrementID()
	shortcut = entities.Shortcut{
		ID:          nextID,
		ShortURL:    shortURL,
		OriginalURL: origURL,
		UserID:      userID,
	}
	if _, err := s.GetOriginalURLByShortcut(ctx, shortURL); err == nil {
		return shortcut, storageerr.ErrNotUniqueShortcut
	} else if _, err = s.GetShortcutByOriginalURL(ctx, origURL); err == nil {
		return shortcut, storageerr.ErrNotUniqueOriginalURL
	}
	s.shortcutListRWMutex.Lock()
	s.shortcutList = append(s.shortcutList, shortcut)
	s.shortcutListRWMutex.Unlock()
	return
}

func (s *Storage) incrementID() uint64 {
	return s.maxID.Add(1)
}

func (s Storage) GetOriginalURLByShortcut(ctx context.Context, shortURL string) (shortcut entities.Shortcut, err error) {
	s.shortcutListRWMutex.RLock()
	defer s.shortcutListRWMutex.RUnlock()
	for _, shortcut = range s.shortcutList {
		if shortcut.ShortURL == shortURL {
			return shortcut, nil
		}
	}
	return shortcut, storageerr.ErrEntityNotFound
}

func (s Storage) GetShortcutByOriginalURL(ctx context.Context, origURL string) (shortcut entities.Shortcut, err error) {
	s.shortcutListRWMutex.RLock()
	defer s.shortcutListRWMutex.RUnlock()
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
	temp := s
	temp.batch = []entities.Shortcut{}
	return temp, nil
}

func (s *Storage) CreateRecordBatch(ctx context.Context, origURL string, shortURL string, userID string) (shortcut entities.Shortcut, err error) {
	shortcut = entities.Shortcut{
		ID:          s.incrementID(),
		UserID:      userID,
		ShortURL:    shortURL,
		OriginalURL: origURL,
		DeletedFlag: false,
	}
	s.batch = append(s.batch, shortcut)
	return
}

func (s *Storage) CommitBatch(ctx context.Context) (err error) {
	s.shortcutListRWMutex.Lock()
	defer s.shortcutListRWMutex.Unlock()
	if err = s.checkIndexesBeforeBatchCommit(); err != nil {
		return err
	}
	for _, sh := range s.batch {
		s.shortcutList = append(s.shortcutList, sh)
	}
	s.batch = nil
	return nil
}

func (s *Storage) RollbackBatch(ctx context.Context) error {
	s.batch = nil
	return nil
}

func (s Storage) GetShortcutsByUser(ctx context.Context, userID string) (shortcuts []entities.Shortcut, err error) {
	s.shortcutListRWMutex.RLock()
	defer s.shortcutListRWMutex.RUnlock()
	for _, sh := range s.shortcutList {
		if sh.UserID == userID {
			shortcuts = append(shortcuts, sh)
		}
	}
	return
}

func (s *Storage) DeleteByShortcutsForUser(ctx context.Context, shortcuts []string, userID string) (err error) {
	s.shortcutListRWMutex.Lock()
	defer s.shortcutListRWMutex.Unlock()
	temp := s.shortcutList
	for idx, sh := range s.shortcutList {
		if sh.UserID == userID && s.isStrSliceContains(shortcuts, sh.ShortURL) {
			temp = s.removeShotcutsListElementByIdx(temp, idx)
		}
	}
	return
}

func (s *Storage) DeleteByShortcutsAndUser(ctx context.Context, deleteTasks []model.DeleteTask) (err error) {
	s.shortcutListRWMutex.Lock()
	defer s.shortcutListRWMutex.Unlock()
	temp := s.shortcutList
	for _, delTask := range deleteTasks {
		for idx, sh := range s.shortcutList {
			if sh.UserID == delTask.UserID && s.isStrSliceContains(delTask.ShortURLs, sh.ShortURL) {
				temp = s.removeShotcutsListElementByIdx(temp, idx)
			}
		}
	}
	s.shortcutList = temp
	return
}

func (s Storage) isStrSliceContains(sl []string, str string) bool {
	for _, a := range sl {
		if a == str {
			return true
		}
	}
	return false
}

func (s Storage) removeShotcutsListElementByIdx(sl []entities.Shortcut, idx int) []entities.Shortcut {
	sl[idx] = sl[len(sl)-1]
	return sl[:len(sl)-1]
}

func (s *Storage) checkIndexesBeforeBatchCommit() (err error) {
	for _, sh := range s.shortcutList {
		for _, batchSh := range s.batch {
			if sh.OriginalURL == batchSh.OriginalURL {
				return storageerr.ErrNotUniqueOriginalURL
			} else if sh.ShortURL == batchSh.ShortURL {
				return storageerr.ErrNotUniqueShortcut
			}
		}
	}
	return
}
