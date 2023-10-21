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
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

const (
	preparedStmtInsertName = "preparedStmtInsertName"
	insertStmt             = `INSERT INTO shortcut (short_url, original_url, user_id) VALUES($1, $2, $3) RETURNING id`
	deleteChunkSize        = 50
	deleteBatchSize        = 100
)

// @todo логика реконекта
type Storage struct {
	conn     *pgxpool.Conn
	connPool *pgxpool.Pool
	logger   *zerolog.Logger
	tx       pgx.Tx
}

func New(ctx context.Context, dsn string, logger *zerolog.Logger) (interfaces.ShortcutStorager, error) {
	logger.Debug().Msg("Connecting to database")

	conPool, err := newConnectionPool(ctx, dsn, logger)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db for storage: %w", err)
	}
	return &Storage{
		connPool: conPool,
		logger:   logger,
	}, nil
}

func (s Storage) NewBatch(ctx context.Context) (interfaces.ShortcutStorager, error) {
	conn, err := s.connPool.Acquire(ctx)
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
		conn:     conn,
		connPool: s.connPool,
		logger:   s.logger,
		tx:       tx,
	}, nil
}

func (s Storage) CommitBatch(ctx context.Context) (err error) {
	if s.tx == nil {
		return errors.New("transaction has not started")
	}
	defer s.conn.Release()
	err = s.tx.Commit(ctx)
	return
}

func (s Storage) RollbackBatch(ctx context.Context) (err error) {
	if s.tx == nil {
		return errors.New("transaction has not started")
	}
	defer s.conn.Release()
	err = s.tx.Rollback(ctx)
	return
}

func (s *Storage) CreateRecordBatch(ctx context.Context, origURL string, shortURL string, userID string) (shortcut entities.Shortcut, err error) {
	var row pgx.Row
	if userID == "" {
		row = s.connPool.QueryRow(ctx, insertStmt, shortURL, origURL, nil)
	} else {
		row = s.connPool.QueryRow(ctx, insertStmt, shortURL, origURL, userID)
	}
	var id uint64
	err = row.Scan(&id)
	if err != nil {
		return shortcut, processInsertStmtError(err)
	}
	return entities.Shortcut{
		ID:          id,
		ShortURL:    shortURL,
		OriginalURL: origURL,
		UserID:      userID,
	}, nil
}

func (s *Storage) CreateRecord(ctx context.Context, origURL string, shortURL string, userID string) (shortcut entities.Shortcut, err error) {
	var row pgx.Row
	if userID == "" {
		row = s.connPool.QueryRow(ctx, insertStmt, shortURL, origURL, nil)
	} else {
		row = s.connPool.QueryRow(ctx, insertStmt, shortURL, origURL, userID)
	}
	var id uint64
	err = row.Scan(&id)
	if err != nil {
		return shortcut, processInsertStmtError(err)
	}
	return entities.Shortcut{
		ID:          id,
		ShortURL:    shortURL,
		OriginalURL: origURL,
		UserID:      userID,
	}, nil
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

func (s Storage) GetOriginalURLByShortcut(ctx context.Context, shortURL string) (shortcut entities.Shortcut, err error) {

	sql := fmt.Sprintf("SELECT * FROM %s WHERE short_url = $1", tableName)
	err = pgxscan.Get(ctx, s.connPool, &shortcut, sql, shortURL)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return shortcut, storageerr.ErrEntityNotFound
	} else if err != nil {
		return shortcut, err
	}
	return
}

func (s Storage) GetShortcutByOriginalURL(ctx context.Context, origURL string) (shortcut entities.Shortcut, err error) {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE original_url = $1", tableName)
	err = pgxscan.Get(ctx, s.connPool, &shortcut, sql, origURL)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return shortcut, storageerr.ErrEntityNotFound
	} else if err != nil {
		return shortcut, err
	}
	return
}

func (s Storage) Close(ctx context.Context) (err error) {
	s.connPool.Close()
	return
}

func (s Storage) Ping(ctx context.Context) (err error) {
	err = s.connPool.Ping(ctx)
	return
}

func (s Storage) GetShortcutsByUser(ctx context.Context, userID string) (shortcuts []entities.Shortcut, err error) {

	sql := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", tableName)
	rows, err := s.connPool.Query(ctx, sql, userID)
	if err != nil {
		return
	}
	defer rows.Close()
	err = pgxscan.ScanAll(&shortcuts, rows)
	return
}

func (s Storage) DeleteByShortcutsForUser(ctx context.Context, shortcuts []string, userID string) (err error) {
	conn, err := s.connPool.Acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()

	sql := fmt.Sprintf("UPDATE %s SET is_deleted = true WHERE user_id = $1 AND short_url = ANY($2)", tableName)

	finished := false
	startIdx := 0
	for {
		batch := &pgx.Batch{}
		for batchIter := 0; batchIter < deleteBatchSize && (batchIter*deleteBatchSize) < len(shortcuts); batchIter++ {
			endIdx := startIdx + deleteChunkSize
			if endIdx > len(shortcuts) {
				endIdx = len(shortcuts)
				finished = true
			}
			batch.Queue(sql, userID, shortcuts[startIdx:endIdx])
			if finished {
				break
			}
			startIdx = endIdx
		}
		batchResults := s.connPool.SendBatch(ctx, batch)
		err = batchResults.Close()
		if err != nil {
			return
		}
		if finished {
			break
		}
	}
	return
}
