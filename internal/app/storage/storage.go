package storage

import (
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/config"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage/filestorage"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage/memstorage"
	"github.com/rs/zerolog"
)

var ShortcutStorage ShortcuterStorage

type ShortcuterStorage interface {
	CreateRecord(origURL string, shortURL string) (err error)
	GetOriginalURLByShortcut(shortURL string) (origURL string, err error)
	GetShortcutByOriginalURL(origURL string) (shortURL string, err error)
	Close() (err error)
}

func ShortcutStorageFactoryInit(opts config.Options, logger *zerolog.Logger) {
	if opts.FileStoragePath == "" {
		ShortcutStorage = memstorage.New()
	} else {
		ShortcutStorage = filestorage.New(opts.FileStoragePath, logger)
	}
}
