package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/manfromth3m0oN/speeddate/pkg/user"
)

// LoginReq is the structure representing a login request
type LoginReq struct {
	Email    string
	Password string
}

// LoginResp is the structure representing a successful login response
type LoginResp struct {
	Token string `json:"token"`
}

// LoginUser performs a login for a user, generating a JWT
func (h *HTTPService) LoginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var loginReq LoginReq
		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		loginInfo, err := user.GetUserLoginInfo(ctx, h.DB, loginReq.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if loginInfo.Password != loginReq.Password {
			http.Error(w, "your password does not match", http.StatusBadRequest)
			return
		}

		// There really should be a refresh token here too
		token, err := jwt.NewBuilder().
			Audience([]string{strconv.Itoa(loginInfo.UserID)}).
			Issuer("speeddate").
			IssuedAt(time.Now()).
			Expiration(time.Now().Add(h.JWTExpr)).
			Build()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		serialized, err := jwt.NewSerializer().Sign(jwt.WithKey(jwa.RS256, h.PrivKey)).Serialize(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := LoginResp{
			Token: string(serialized),
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
