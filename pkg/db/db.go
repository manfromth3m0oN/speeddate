package db

import (
	"database/sql"

	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/lib/pq"
	"github.com/manfromth3m0oN/speeddate/cmd/config"
)

func Connect(cfg config.Config) (*goqu.Database, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/postgres?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	gdb := goqu.New("postgres", db)
	return gdb, nil
}
