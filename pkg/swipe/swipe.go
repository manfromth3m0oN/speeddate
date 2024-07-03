package swipe

import (
	"context"
	"log/slog"

	"github.com/doug-martin/goqu/v9"
)

type Swipe struct {
	Id         int
	SwiperId   int
	SwipeeId   int
	Preference bool
}

func (s *Swipe) Insert(ctx context.Context, db *goqu.Database) error {
	insertSQL, _, err := goqu.Insert("swipe").Rows(
		goqu.Record{
			"id":         goqu.Default(),
			"swiper_id":  s.SwiperId,
			"swipee_id":  s.SwipeeId,
			"preference": s.Preference,
			"insert_ts":  goqu.Default(),
		}).Returning("id").ToSQL()
	if err != nil {
		return err
	}

	results, err := db.QueryContext(ctx, insertSQL)
	if err != nil {
		return err
	}

	for results.Next() {
		if err := results.Scan(&s.Id); err != nil {
			return err
		}
	}

	return nil
}

func FindReciprocation(ctx context.Context, db *goqu.Database, swiper, swipee int) (bool, error) {
	getSQL, _, err := goqu.From("swipe").
		Select("preference").
		Where(goqu.C("swiper_id").Eq(swipee), goqu.C("swipee_id").Eq(swiper)).ToSQL()
	if err != nil {
		slog.Error("failed to build query sql", "err", err)
		return false, err
	}

	results, err := db.QueryContext(ctx, getSQL)
	if err != nil {
		slog.Error("failed to perform reciprocation query", "err", err)
		return false, err
	}

	var preference bool
	for results.Next() {
		if err := results.Scan(&preference); err != nil {
			slog.Error("failed to scan query result", "err", err)
			return false, nil
		}
	}

	return preference, nil
}
