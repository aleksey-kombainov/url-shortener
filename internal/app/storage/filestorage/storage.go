package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/entities"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/interfaces"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage/storageerr"
	"github.com/rs/zerolog"
	"os"
)

const (
	recordsSeparator = '\n'
)

type Storage struct {
	shortcutList []entities.Shortcut
	maxID        uint64
	fileHdl      *os.File
	logger       *zerolog.Logger
}

func New(fileStoragePath string, logger *zerolog.Logger) *Storage {
	logger.Debug().Msgf("Opening '%s' storage", fileStoragePath)
	fileHdl, err := os.OpenFile(fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.Fatal().Msg("Open file error. " + err.Error())
	}
	s := &Storage{
		shortcutList: make([]entities.Shortcut, 0),
		maxID:        0,
		fileHdl:      fileHdl,
		logger:       logger,
	}
	s.loadData()
	return s
}

func (s Storage) NewBatcher(ctx context.Context) (*Storage, error) {
	return &s, nil
}

func (s *Storage) CreateRecord(ctx context.Context, origURL string, shortURL string) (err error) {
	s.maxID++
	rec := entities.Shortcut{
		ID:          s.maxID,
		ShortURL:    shortURL,
		OriginalURL: origURL,
	}
	if _, err = s.GetOriginalURLByShortcut(ctx, shortURL); err == nil {
		return storageerr.ErrNotUniqueShortcut
	} else if _, err = s.GetShortcutByOriginalURL(ctx, origURL); err == nil {
		return storageerr.ErrNotUniqueOriginalURL
	}
	dataStr, err := json.Marshal(rec)
	if err != nil {
		s.logger.Error().Msg("can not marshal new record: " + err.Error())
		return
	}
	dataStr = append(dataStr, recordsSeparator)
	if _, err = s.fileHdl.Write(dataStr); err != nil {
		s.logger.Error().Msg("can not write new record: " + err.Error())
		return
	}
	s.shortcutList = append(s.shortcutList, rec)

	return nil
}

func (s Storage) GetOriginalURLByShortcut(ctx context.Context, shortURL string) (origURL string, err error) {
	for _, sh := range s.shortcutList {
		if sh.ShortURL == shortURL {
			return sh.OriginalURL, nil
		}
	}
	return "", storageerr.ErrEntityNotFound
}

func (s Storage) GetShortcutByOriginalURL(tx context.Context, origURL string) (shortURL string, err error) {
	for _, sh := range s.shortcutList {
		if sh.OriginalURL == origURL {
			return sh.ShortURL, nil
		}
	}
	return "", storageerr.ErrEntityNotFound
}

func (s *Storage) Close(tx context.Context) (err error) {
	s.logger.Info().Msg("closing storage file")
	if err = s.fileHdl.Close(); err != nil {
		s.logger.Error().Msg("Cant close storage: " + err.Error())
		return
	}
	return nil
}

func (s Storage) Ping(ctx context.Context) (err error) {
	return nil
}

func (s *Storage) loadData() {
	scanner := bufio.NewScanner(s.fileHdl)
	// одиночное сканирование до следующей строки
	for scanner.Scan() {
		sc := entities.Shortcut{}
		bytesSlice := scanner.Bytes()
		err := json.Unmarshal(bytesSlice, &sc)
		if err != nil {
			s.logger.Fatal().Msg("unmarshall error while scanning file: " + err.Error())
		}
		if sc.ID > s.maxID {
			s.maxID = sc.ID
		}
		s.shortcutList = append(s.shortcutList, sc)
	}
	if err := scanner.Err(); err != nil {
		s.logger.Fatal().Msg("Error while scanning file: " + err.Error())
	}
	s.logger.Info().Msgf("Loaded %d records from storage", len(s.shortcutList))
}

func (s Storage) NewBatch(ctx context.Context) (interfaces.ShortcutStorager, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) CreateRecordBatch(ctx context.Context, origURL string, shortURL string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) CommitBatch(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) RollbackBatch(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
