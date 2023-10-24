package interfaces

import (
	"context"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/entities"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/model"
)

type ShortcutStorager interface {
	CreateRecord(ctx context.Context, origURL string, shortURL string, userID string) (shortcut entities.Shortcut, err error)

	NewBatch(ctx context.Context) (ShortcutStorager, error)
	CreateRecordBatch(ctx context.Context, origURL string, shortURL string, userID string) (shortcut entities.Shortcut, err error)
	CommitBatch(ctx context.Context) error
	RollbackBatch(ctx context.Context) error

	GetOriginalURLByShortcut(ctx context.Context, shortURL string) (shortcut entities.Shortcut, err error)
	GetShortcutByOriginalURL(ctx context.Context, origURL string) (shortcut entities.Shortcut, err error)
	GetShortcutsByUser(ctx context.Context, userID string) (shortcuts []entities.Shortcut, err error)
	Close(ctx context.Context) (err error)
	Ping(ctx context.Context) (err error)

	DeleteByShortcutsForUser(ctx context.Context, shortcuts []string, userID string) (err error)

	DeleteByShortcutsAndUser(ctx context.Context, deleteTasks []model.DeleteTask) (err error)
}
