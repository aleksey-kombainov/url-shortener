package storage

import (
	"context"
	"errors"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage/filestorage"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage/memstorage"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage/postgres"
	"github.com/rs/zerolog"
)

const (
	TypeDB     = "db"
	TypeFile   = "file"
	TypeMemory = "memory"
)

var (
	EntityNotFoundErr = errors.New("Entity not found")
)

type ShortcutStorager interface {
	CreateRecord(origURL string, shortURL string) (err error)
	GetOriginalURLByShortcut(shortURL string) (origURL string, err error)
	GetShortcutByOriginalURL(origURL string) (shortURL string, err error)
	Close() (err error)
	Ping(ctx context.Context) (err error)
}

func ShortcutStorageFactory(ctx context.Context, logger *zerolog.Logger, storageType string, param string) (storage ShortcutStorager, err error) {
	switch storageType {
	case TypeDB:
		storage, err = postgres.New(ctx, param, logger, EntityNotFoundErr)
	case TypeFile:
		storage = filestorage.New(param, logger, EntityNotFoundErr)
	case TypeMemory:
		storage = memstorage.New(EntityNotFoundErr)
	default:
		err = errors.New("unknown storage type")
	}
	return
}
