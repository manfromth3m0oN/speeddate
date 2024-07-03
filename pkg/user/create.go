package user

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

// Insert inserts a user into the database
func (u *User) Insert(ctx context.Context, db *goqu.Database) error {
	insertDs := goqu.Insert("users").Rows(
		goqu.Record{
			"id":        goqu.Default(),
			"email":     u.Email,
			"password":  u.Password,
			"name":      u.Name,
			"gender":    u.Gender,
			"age":       u.Age,
			"longitude": u.longitude,
			"latitude":  u.latitude,
			"insert_ts": goqu.Default(),
		}).Returning("id")

	insertSQL, _, err := insertDs.ToSQL()
	if err != nil {
		return err
	}

	result, err := db.QueryContext(ctx, insertSQL)
	if err != nil {
		return err
	}

	for result.Next() {
		if err := result.Scan(&u.Id); err != nil {
			return err
		}
	}

	return nil
}
