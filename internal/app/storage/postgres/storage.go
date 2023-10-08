package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/entities"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

const tableName = "shortcut"

type Storage struct {
	conn              *pgx.Conn
	logger            *zerolog.Logger
	entityNotFoundErr error
}

func New(ctx context.Context, dsn string, logger *zerolog.Logger, entityNotFoundErr error) (*Storage, error) {
	logger.Debug().Msg("Connecting to database")

	conn, err := NewConnection(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db for storage: %w", err)
	}

	return &Storage{
		conn:              conn,
		logger:            logger,
		entityNotFoundErr: entityNotFoundErr,
	}, nil
}

func (s *Storage) CreateRecord(origURL string, shortURL string) (err error) {
	sql := fmt.Sprintf("INSERT INTO %s (short_url, original_url) VALUES($1, $2)", tableName)
	_, err = s.conn.Exec(context.Background(), sql, shortURL, origURL)
	if err != nil {
		return fmt.Errorf("can't write record to db for storage: %w", err)
	}
	return
}

func (s Storage) GetOriginalURLByShortcut(shortURL string) (origURL string, err error) {
	sql := fmt.Sprintf("SELECT id, short_url, original_url FROM %s WHERE short_url = $1", tableName)
	var entity entities.Shortcut
	err = pgxscan.Get(context.Background(), s.conn, &entity, sql, shortURL)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", s.entityNotFoundErr
	} else if err != nil {
		return "", err
	}
	return entity.OriginalURL, nil
}

func (s Storage) GetShortcutByOriginalURL(origURL string) (shortURL string, err error) {
	sql := fmt.Sprintf("SELECT id, short_url, original_url FROM %s WHERE original_url = $1", tableName)
	var entity entities.Shortcut
	err = pgxscan.Get(context.Background(), s.conn, &entity, sql, origURL)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", s.entityNotFoundErr
	} else if err != nil {
		return "", err
	}
	return entity.ShortURL, nil
}

func (s *Storage) Close() (err error) {
	err = s.conn.Close(context.Background())
	return
}

func (s *Storage) Ping(ctx context.Context) (err error) {
	return s.conn.Ping(ctx)
}
