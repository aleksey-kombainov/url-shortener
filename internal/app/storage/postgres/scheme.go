package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

const (
	tableName              = "shortcut"
	shortcutIdxShortURL    = "shortcut__index_short_url"
	shortcutIdxOriginalURL = "shortcut__index_original_url"
	shortcutIdxUserID      = "shortcut__index_user_id"
)

func createScheme(ctx context.Context, conn *pgxpool.Conn, logger *zerolog.Logger) error {
	dqlQueries := []string{
		fmt.Sprintf(`create table %s
			(
				id           bigint primary key generated always as identity,
				user_id uuid,
				short_url    char(8) not null,
				original_url varchar(255) not null,
    			is_deleted BOOLEAN NOT NULL DEFAULT FALSE
			)`, tableName),
		fmt.Sprintf("create unique index %s	on shortcut (short_url)", shortcutIdxShortURL),
		fmt.Sprintf("create unique index %s on shortcut (original_url)", shortcutIdxOriginalURL),
		fmt.Sprintf("create index %s on shortcut (user_id)", shortcutIdxUserID),
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("can't start transaction: %w", err)
	}
	for _, queryStr := range dqlQueries {
		if _, err := tx.Exec(ctx, queryStr); err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return fmt.Errorf("can't execute query\n%s \nand even rollback it: %w", queryStr, err)
			}
			return fmt.Errorf("can't execute query\n%s: %w", queryStr, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("can't commit transaction: %w", err)
	}
	logger.Info().Msg("db migrations successfully applied!")
	return nil
}
