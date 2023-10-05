package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

func GetConnection(ctx context.Context, dsn string, logger zerolog.Logger) (conn *pgx.Conn, err error) {
	conn, err = pgx.Connect(ctx, dsn)
	if err != nil {
		logger.Error().Msgf("Unable to connect to database: %v\n", err)
	}
	return
}
