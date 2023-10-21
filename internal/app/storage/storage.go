package storage

import (
	"context"
	"errors"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/interfaces"
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

func ShortcutStorageFactory(ctx context.Context, logger *zerolog.Logger, storageType string, param string) (storage interfaces.ShortcutStorager, err error) {
	switch storageType {
	case TypeDB:
		storage, err = postgres.New(ctx, param, logger)
	case TypeFile:
		storage = filestorage.New(param, logger)
	case TypeMemory:
		storage = memstorage.New()
	default:
		err = errors.New("unknown storage type")
	}
	return
}
