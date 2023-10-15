package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

const (
	tableName              = "shortcut"
	shortcutIdxShortURL    = "shortcut__index_short_url"
	shortcutIdxOriginalURL = "shortcut__index_original_url"
)

func createScheme(ctx context.Context, conn *pgx.Conn, logger *zerolog.Logger) error {
	dqlQueries := []string{
		fmt.Sprintf(`create table %s
			(
				id           bigint primary key generated always as identity,
				short_url    char(8) not null,
				original_url varchar(255) not null
			)`, tableName),
		fmt.Sprintf("create unique index %s	on shortcut (short_url)", shortcutIdxShortURL),
		fmt.Sprintf("create unique index %s on shortcut (original_url)", shortcutIdxOriginalURL),
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
