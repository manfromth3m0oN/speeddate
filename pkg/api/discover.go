package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"github.com/manfromth3m0oN/speeddate/pkg/user"
)

// DiscoverUsers finds all users but the currently logged in user
// There are filters for age (age_high and age_low) as well as gender
// Distances from the current user are also calculated
func (h *HTTPService) DiscoverUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userId, ok := ctx.Value(UserIDKey).(int)
		if !ok {
			http.Error(w, "couldnt get user id, are you logged in?", http.StatusInternalServerError)
			return
		}

		queryParams := r.URL.Query()
		filters, err := constructFilters(queryParams)
		if err != nil {
			return
		}

		users, err := user.GetAllOtherUsers(ctx, h.DB, userId, filters...)
		if err != nil {
			slog.Info("failed fetching other users", "err", err)
			http.Error(w, "failed fetching other users", http.StatusInternalServerError)
			return
		}

		userLat, userLong, err := user.GetUserCoords(ctx, h.DB, userId)
		if err != nil {
			slog.Info("failed to get coordinates for loggedin user", "err", err)
			http.Error(w, "failed to get coordinates for loggedin user", http.StatusInternalServerError)
			return
		}

		for i := 0; i < len(users); i++ {
			users[i].CalculateDistance(userLat, userLong)
		}

		sort.Sort(user.ByDistance(users))

		slog.Info("users", "users", users)

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "failed to encode users json", http.StatusInternalServerError)
			return
		}
	}
}

// constructFilters for gender, age_high and age_low query params, construct goqu filters for the database
func constructFilters(queryParams url.Values) ([]user.OtherUsersFilter, error) {
	ageLowFilter := queryParams.Get("age_low")
	ageHighFilter := queryParams.Get("age_high")
	genderFilter := queryParams.Get("gender")

	filters := make([]user.OtherUsersFilter, 0)
	if ageLowFilter != "" && ageHighFilter != "" {
		low, err := strconv.Atoi(ageLowFilter)
		if err != nil {
			return nil, err
		}

		high, err := strconv.Atoi(ageHighFilter)
		if err != nil {
			return nil, err
		}

		filters = append(filters, user.RangeFilter{Attr: "age", Low: low, High: high})
	}

	if genderFilter != "" {
		filters = append(filters, user.EqFilter{Attr: "gender", Value: genderFilter})
	}

	slog.Info("filters", "filters", filters)
	return filters, nil
}
