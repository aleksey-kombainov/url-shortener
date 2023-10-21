package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func newConnectionPool(ctx context.Context, dsn string, logger *zerolog.Logger) (conPool *pgxpool.Pool, err error) {
	//if strings.Contains(dsn, "?") {
	//	dsn = dsn + "&"
	//} else {
	//	dsn = dsn + "?"
	//}
	//dsn += "&pool_max_conns=50"
	conPool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db for storage: %w", err)
	}
	con, err := conPool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	if err = con.Ping(ctx); err != nil {
		return nil, fmt.Errorf("connection to database established, but ping() returned error: %w", err)
	}

	err = checkDBSetup(ctx, con, logger)
	if err != nil {
		return
	}

	con.Release()
	return
}

func checkDBSetup(ctx context.Context, conn *pgxpool.Conn, logger *zerolog.Logger) (err error) {
	exists, err := tableExists(ctx, conn, "shortcut")
	if err != nil {
		return
	}
	if exists {
		return nil
	}
	err = createScheme(ctx, conn, logger)
	return
}

func tableExists(ctx context.Context, conn *pgxpool.Conn, tableName string) (bool, error) {
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
