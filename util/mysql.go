package util

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

func OpenMySQLWithMultiStatements(dsn string) (*sql.DB, error) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	cfg.MultiStatements = true
	return sql.Open("mysql", cfg.FormatDSN())
}

func MigrateMySQL(ctx context.Context, dsn string, stmt string) error {
	db, err := OpenMySQLWithMultiStatements(dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.ExecContext(ctx, stmt)
	return err
}
