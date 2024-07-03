package match

import (
	"context"

	"github.com/doug-martin/goqu/v9"
)

type Match struct {
	Id    int
	UserA int
	UserB int
}

// Insert inserts a match into the DB
func (m *Match) Insert(ctx context.Context, db *goqu.Database) error {
	insertSQL, _, err := goqu.
		Insert("match").
		Rows(goqu.Record{
			"id":        goqu.Default(),
			"user_a":    m.UserA,
			"user_b":    m.UserB,
			"insert_ts": goqu.Default(),
		}).
		Returning("id").ToSQL()
	if err != nil {
		return err
	}

	result, err := db.QueryContext(ctx, insertSQL)
	if err != nil {
		return err
	}

	for result.Next() {
		if err := result.Scan(&m.Id); err != nil {
			return err
		}
	}

	return nil
}
