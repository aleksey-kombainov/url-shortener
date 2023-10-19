package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/entities"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/interfaces"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage/storageerr"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
)

const (
	preparedStmtInsertName = "preparedStmtInsertName"
	insertStmt             = `INSERT INTO shortcut (short_url, original_url, user_id) VALUES($1, $2, $3) RETURNING id`
)

// @todo логика реконекта
type Storage struct {
	dsn    string
	conn   *pgx.Conn
	logger *zerolog.Logger
	tx     pgx.Tx
}

func New(ctx context.Context, dsn string, logger *zerolog.Logger) (interfaces.ShortcutStorager, error) {
	logger.Debug().Msg("Connecting to database")

	conn, err := newConnection(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db for storage: %w", err)
	}

	if err = checkDBSetup(ctx, conn, logger); err != nil {
		return nil, err
	}

	return &Storage{
		dsn:    dsn,
		conn:   conn,
		logger: logger,
	}, nil
}

func (s Storage) NewBatch(ctx context.Context) (interfaces.ShortcutStorager, error) {
	conn, err := newConnection(ctx, s.dsn)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db for storage: %w", err)
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %w", err)
	}
	_, err = tx.Prepare(ctx, preparedStmtInsertName, insertStmt)
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %w", err)
	}
	return &Storage{
		dsn:    s.dsn,
		conn:   conn,
		logger: s.logger,
		tx:     tx,
	}, nil
}

func (s Storage) CommitBatch(ctx context.Context) (err error) {
	if s.tx == nil {
		return errors.New("transaction has not started")
	}
	err = s.tx.Commit(ctx)
	return
}

func (s Storage) RollbackBatch(ctx context.Context) (err error) {
	if s.tx == nil {
		return errors.New("transaction has not started")
	}
	err = s.tx.Rollback(ctx)
	return
}

func (s *Storage) CreateRecordBatch(ctx context.Context, origURL string, shortURL string, userID string) (err error) {
	if userID == "" {
		_, err = s.tx.Exec(ctx, preparedStmtInsertName, shortURL, origURL, nil)
	} else {
		_, err = s.tx.Exec(ctx, preparedStmtInsertName, shortURL, origURL, userID)
	}
	if err != nil {
		return processInsertStmtError(err)
	}
	return
}

func (s *Storage) CreateRecord(ctx context.Context, origURL string, shortURL string, userID string) (err error) {
	if userID == "" {
		_, err = s.conn.Exec(ctx, insertStmt, shortURL, origURL, nil)
	} else {
		_, err = s.conn.Exec(ctx, insertStmt, shortURL, origURL, userID)
	}
	if err != nil {
		return processInsertStmtError(err)
	}
	return
}

func processInsertStmtError(err error) error {
	var pgErr *pgconn.PgError
	isPgErr := errors.As(err, &pgErr)
	if isPgErr {
		if pgErr.ConstraintName == shortcutIdxShortURL {
			return storageerr.ErrNotUniqueShortcut
		} else if pgErr.ConstraintName == shortcutIdxOriginalURL {
			return storageerr.ErrNotUniqueOriginalURL
		} else {
			return fmt.Errorf("can't write record to db for storage: %w", pgErr)
		}
	} else {
		return fmt.Errorf("can't write record to db for storage: %w", err)
	}
}

func (s Storage) GetOriginalURLByShortcut(ctx context.Context, shortURL string) (origURL string, err error) {
	sql := fmt.Sprintf("SELECT id, short_url, original_url FROM %s WHERE short_url = $1", tableName)
	var entity entities.Shortcut
	err = pgxscan.Get(ctx, s.conn, &entity, sql, shortURL)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", storageerr.ErrEntityNotFound
	} else if err != nil {
		return "", err
	}
	return entity.OriginalURL, nil
}

func (s Storage) GetShortcutByOriginalURL(ctx context.Context, origURL string) (shortURL string, err error) {
	sql := fmt.Sprintf("SELECT id, short_url, original_url FROM %s WHERE original_url = $1", tableName)
	var entity entities.Shortcut
	err = pgxscan.Get(ctx, s.conn, &entity, sql, origURL)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", storageerr.ErrEntityNotFound
	} else if err != nil {
		return "", err
	}
	return entity.ShortURL, nil
}

func (s Storage) Close(ctx context.Context) (err error) {
	err = s.conn.Close(ctx)
	return
}

func (s Storage) Ping(ctx context.Context) (err error) {
	return s.conn.Ping(ctx)
}

func (s Storage) GetShortcutsByUser(ctx context.Context, userID string) (shortcuts []entities.Shortcut, err error) {
	sql := fmt.Sprintf("SELECT id, short_url, original_url FROM %s WHERE user_id = $1", tableName)
	rows, err := s.conn.Query(ctx, sql, userID)
	if err != nil {
		return
	}
	err = pgxscan.ScanAll(&shortcuts, rows)
	return
}
