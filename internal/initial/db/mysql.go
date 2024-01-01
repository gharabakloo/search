package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // nolint: revive

	"gharabakloo/search/internal/entity"
	"gharabakloo/search/pkg/myerr"
)

func MySQL(ctx context.Context, cfg entity.MySQLConfig) (*sql.DB, error) {
	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DBName)
	db, err := sql.Open(cfg.Driver, sourceName)
	if err != nil {
		return nil, myerr.Errorf(err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, myerr.Errorf(err)
	}
	return db, nil
}
