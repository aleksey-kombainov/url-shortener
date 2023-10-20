package app

import (
	"context"
	"fmt"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/entities"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/interfaces"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/model"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/random"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage/storageerr"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"time"
)

const (
	generatorIterationLimit = 1000
	shortcutLength          = 8
)

type storageSaver func(ctx context.Context, origURL string, shortURL string, userID string) (shortcut entities.Shortcut, err error)

type ShortcutService struct {
	logger  *zerolog.Logger
	Storage *interfaces.ShortcutStorager
}

func NewShortcutService(logger *zerolog.Logger, storage interfaces.ShortcutStorager) *ShortcutService {
	return &ShortcutService{
		logger:  logger,
		Storage: &storage,
	}
}

func (s ShortcutService) MakeShortcut(url string, userID string) (shortcut entities.Shortcut, err error) {
	shortcut, err = s.generateAndSaveShortcut(url, userID, (*s.Storage).CreateRecord)
	if errors.Is(err, storageerr.ErrNotUniqueOriginalURL) {
		shortcut, errGettingShortcut := (*s.Storage).GetShortcutByOriginalURL(context.TODO(), url)
		if errGettingShortcut != nil {
			return shortcut, errGettingShortcut
		}
		return shortcut, err
	}
	return
}

func (s ShortcutService) generateAndSaveShortcut(url string, userID string, saveMethod storageSaver) (shortcut entities.Shortcut, err error) {
	if url == "" {
		return shortcut, errors.New("url is empty")
	}
	isGenerated := false
	for i := 0; i < generatorIterationLimit; i++ {
		sh := random.GenString(shortcutLength)

		shortcut, err = saveMethod(context.TODO(), url, sh, userID)

		if err == nil {
			isGenerated = true
			break
		} else if errors.Is(err, storageerr.ErrNotUniqueShortcut) {
			continue
		} else {
			//s.logger.Error().Msgf("creating shortcut - error while creating shortcut: %w", err)
			return shortcut, fmt.Errorf("creating shortcut - error while creating shortcut: %w", err)
		}
	}
	if !isGenerated {
		return shortcut, errors.New("generator limit exceeded")
	}
	return
}

func (s ShortcutService) MakeShortcutBatch(ctx context.Context, batch []model.ShortenerBatchRecordRequest, userID string) (result []model.ShortenerBatchRecordResponse, err error) {
	batchStorage, err := (*s.Storage).NewBatch(ctx)
	if err != nil {
		return
	}
	defer func() {
		if err := batchStorage.Close(ctx); err != nil {
			s.logger.Error().Msgf("can't close storage connection: %s", err)
		}
	}()
	for _, batchRecord := range batch {
		shortcut, err := s.generateAndSaveShortcut(batchRecord.OriginalURL, userID, batchStorage.CreateRecordBatch)
		if err != nil {
			if errRollback := batchStorage.RollbackBatch(context.TODO()); errRollback != nil {
				return nil, fmt.Errorf("%w; can't rollback transaction: %w", err, errRollback)
			}
			return nil, err
		}
		result = append(result, model.ShortenerBatchRecordResponse{CorrelationID: batchRecord.CorrelationID, ShortURL: shortcut.ShortURL})
	}
	err = batchStorage.CommitBatch(context.TODO())
	return
}

func (s ShortcutService) GetOriginalURLByShortcut(shortURL string) (shortcut entities.Shortcut, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return (*s.Storage).GetOriginalURLByShortcut(ctx, shortURL)
}

func (s ShortcutService) GetShortcutByOriginalURL(origURL string) (shortcut entities.Shortcut, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return (*s.Storage).GetShortcutByOriginalURL(ctx, origURL)
}

func (s ShortcutService) GetShortcutsByUser(userID string) (shortcuts []entities.Shortcut, err error) {
	shortcuts, err = (*s.Storage).GetShortcutsByUser(context.Background(), userID)
	return
}

func (s ShortcutService) DeleteByIDsAndUser(shortcuts []string, userID string) (err error) {
	if userID == "" {
		return fmt.Errorf("invalid userID")
	}
	err = (*s.Storage).DeleteByShortcutsForUser(context.Background(), shortcuts, userID)
	return
}
