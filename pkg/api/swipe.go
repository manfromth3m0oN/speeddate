package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/manfromth3m0oN/speeddate/pkg/match"
	"github.com/manfromth3m0oN/speeddate/pkg/swipe"
)

// SwipeReq represents a swipe request
type SwipeReq struct {
	Swipee     int
	Preference bool
}

// SwipeResp represents a swipe response
type SwipeResp struct {
	Matched bool `json:"matched"`
	MatchId int  `json:"matchID,omitempty"`
}

// SwipeOnUser is a HTTP handler that performs a swipe on user
// for the user the AuthMiddleware identifies
func (h *HTTPService) SwipeOnUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userId, ok := ctx.Value(UserIDKey).(int)
		if !ok {
			slog.Error("couldnt get user id")
			http.Error(w, "couldnt get user id, are you logged in?", http.StatusInternalServerError)
			return
		}

		var swipeReq SwipeReq
		if err := json.NewDecoder(r.Body).Decode(&swipeReq); err != nil {
			slog.Error("failed to decode swipe request", "err", err)
			http.Error(w, "failed to decode swipe request", http.StatusBadRequest)
			return
		}

		s := swipe.Swipe{
			SwiperId:   userId,
			SwipeeId:   swipeReq.Swipee,
			Preference: swipeReq.Preference,
		}
		if err := s.Insert(ctx, h.DB); err != nil {
			slog.Error("failed to insert swipe", "err", err)
			http.Error(w, "failed to insert swipe", http.StatusInternalServerError)
			return
		}

		switch s.Preference {
		case true:
			slog.Info("user preference is true")
			matched, err := swipe.FindReciprocation(ctx, h.DB, s.SwiperId, s.SwipeeId)
			if err != nil {
				slog.Error("failed in finding reciprocations", "err", err)
				http.Error(w, "failed in finding reciprocations", http.StatusInternalServerError)
			}
			if matched {
				m := match.Match{
					UserA: s.SwiperId,
					UserB: s.SwipeeId,
				}
				if err := m.Insert(ctx, h.DB); err != nil {
					slog.Error("failed to create match", "err", err)
					http.Error(w, "failed to create match", http.StatusInternalServerError)
					return
				}
				if err := json.NewEncoder(w).Encode(SwipeResp{Matched: true, MatchId: m.Id}); err != nil {
					slog.Error("failed to encode no match response", "err", err)
					http.Error(w, "failed to encode no match response", http.StatusInternalServerError)
					return
				}
			}
		case false:
			slog.Info("user preference is true")
			if err := json.NewEncoder(w).Encode(SwipeResp{Matched: false}); err != nil {
				slog.Error("failed to encode no match response", "err", err)
				http.Error(w, "failed to encode no match response", http.StatusInternalServerError)
				return
			}
		}
	}
}
