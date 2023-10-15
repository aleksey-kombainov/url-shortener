package interfaces

import "context"

type ShortcutStorager interface {
	CreateRecord(ctx context.Context, origURL string, shortURL string) (err error)

	NewBatch(ctx context.Context) (ShortcutStorager, error)
	CreateRecordBatch(ctx context.Context, origURL string, shortURL string) (err error)
	CommitBatch(ctx context.Context) error
	RollbackBatch(ctx context.Context) error

	GetOriginalURLByShortcut(ctx context.Context, shortURL string) (origURL string, err error)
	GetShortcutByOriginalURL(ctx context.Context, origURL string) (shortURL string, err error)
	Close(ctx context.Context) (err error)
	Ping(ctx context.Context) (err error)
}
