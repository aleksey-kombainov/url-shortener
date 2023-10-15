package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

func newConnection(ctx context.Context, dsn string) (conn *pgx.Conn, err error) {
	conn, err = pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	if err = conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("connection to database established, but ping() returned error: %w", err)
	}
	return conn, nil
}

func checkDBSetup(ctx context.Context, conn *pgx.Conn, logger *zerolog.Logger) (err error) {
	exists, err := tableExists(ctx, conn, "shortcut")
	if exists {
		return nil
	}
	err = createScheme(ctx, conn, logger)
	return
}

func tableExists(ctx context.Context, conn *pgx.Conn, tableName string) (bool, error) {
	sql := `SELECT EXISTS (
		SELECT FROM 
			information_schema.tables 
		WHERE 
			table_schema LIKE 'public' AND 
			table_type LIKE 'BASE TABLE' AND
			table_name = $1
    )`
	row := conn.QueryRow(ctx, sql, tableName)
	var res bool
	err := row.Scan(&res)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return res, nil
}
