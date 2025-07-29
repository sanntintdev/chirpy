package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sanntintdev/chirpy/internal/auth"
	"github.com/sanntintdev/chirpy/internal/database"
	"github.com/sanntintdev/chirpy/pkg/config"
	"github.com/sanntintdev/chirpy/pkg/middleware"
)

type AccessTokenResponse struct {
	Token string `json:"token"`
}

func RotateTokenHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		refreshTokenData, ok := ctx.Value(middleware.RefreshTokenKey).(database.GetRefreshTokenRow)
		if !ok {
			http.Error(w, "Invalid refresh token context", http.StatusInternalServerError)
			return
		}

		user, err := cfg.DBQueries.GetUserByID(ctx, refreshTokenData.UserID)
		if err != nil {
			http.Error(w, "User not found", http.StatusInternalServerError)
			return
		}

		tokenService := auth.NewTokenService([]byte(cfg.JWT_SECRET_KEY), cfg.JWT_ISSUER)
		accessToken, err := tokenService.GenerateToken(user)
		if err != nil {
			http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
			return
		}

		response := AccessTokenResponse{
			Token: accessToken,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func RevokeTokenHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		refreshTokenData, ok := ctx.Value(middleware.RefreshTokenKey).(database.GetRefreshTokenRow)
		if !ok {
			http.Error(w, "Invalid refresh token context", http.StatusInternalServerError)
			return
		}

		err := cfg.DBQueries.RevokeRefreshToken(ctx, refreshTokenData.Token)
		if err != nil {
			http.Error(w, "Failed to revoke token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
