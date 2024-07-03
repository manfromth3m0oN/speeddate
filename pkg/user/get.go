package user

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

// LoginInfo holds the UserID and Password for a user, for comparison during login
type LoginInfo struct {
	UserID   int
	Password string
}

// GetUserLoginInfo retrives a users login info fro mthe database for a given username
func GetUserLoginInfo(ctx context.Context, db *goqu.Database, username string) (LoginInfo, error) {
	getSQL, _, err := goqu.From("users").Select("password", "id").Where(goqu.C("email").Eq(username)).ToSQL()
	if err != nil {
		return LoginInfo{}, err
	}

	result, err := db.QueryContext(ctx, getSQL)
	if err != nil {
		return LoginInfo{}, err
	}
	defer result.Close()

	loginInfo := LoginInfo{}
	for result.Next() {
		if err := result.Scan(&loginInfo.Password, &loginInfo.UserID); err != nil {
			return LoginInfo{}, err
		}
	}

	return loginInfo, nil
}

// OtherUsersFilter is an interface that any filter can implement
// This means you can have an artibitrarily complex filter that turns into a goqu expression for easy querying
type OtherUsersFilter interface {
	ToExpression() exp.Expression
}

// EqFilter is a filter for equality
type EqFilter struct {
	Attr  string
	Value any
}

func (e EqFilter) ToExpression() exp.Expression {
	return goqu.C(e.Attr).Eq(e.Value)
}

// RangeFilter is a filter for values in a range
type RangeFilter struct {
	Attr string
	Low  int
	High int
}

func (r RangeFilter) ToExpression() exp.Expression {
	return goqu.C(r.Attr).Between(exp.NewRangeVal(r.Low, r.High))
}

// GetAllOtherUsers gets all users but the one supplied in currentUserId
func GetAllOtherUsers(ctx context.Context, db *goqu.Database, currentUserId int, filters ...OtherUsersFilter) ([]User, error) {
	whereExpressions := []exp.Expression{goqu.C("id").Neq(currentUserId)}
	for _, filter := range filters {
		whereExpressions = append(whereExpressions, filter.ToExpression())
	}

	getSQL, _, err := goqu.From("users").Select("id", "name", "gender", "age", "latitude", "longitude").Where(whereExpressions...).ToSQL()
	if err != nil {
		return nil, err
	}

	result, err := db.QueryContext(ctx, getSQL)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	users := make([]User, 0)
	for result.Next() {
		var user User
		if err := result.Scan(&user.Id, &user.Name, &user.Gender, &user.Age, &user.latitude, &user.longitude); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserCoords gets the users coordinates from the database
func GetUserCoords(ctx context.Context, db *goqu.Database, userId int) (float64, float64, error) {
	getSQL, _, err := goqu.From("users").Select("latitude", "longitude").Where(goqu.C("id").Eq(userId)).ToSQL()
	if err != nil {
		return 0, 0, err
	}

	result, err := db.QueryContext(ctx, getSQL)
	if err != nil {
		return 0, 0, err
	}

	var lat, long float64
	for result.Next() {
		if err := result.Scan(&lat, &long); err != nil {
			return 0, 0, err
		}
	}

	return lat, long, nil
}
