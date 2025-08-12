package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/sanntintdev/chirpy/internal/auth"
	"github.com/sanntintdev/chirpy/pkg/config"
)

func MetricsInc(cfg *config.APIConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHitCount.Add(1)
		next.ServeHTTP(w, r)
	})
}

type contextKey string

const UserIDKey contextKey = "user_id"
const RefreshTokenKey contextKey = "refresh_token"

func AuthMiddleware(cfg *config.APIConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenService := auth.NewTokenService([]byte(cfg.JWT_SECRET_KEY), cfg.JWT_ISSUER)

		userID, err := tokenService.ValidateToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateRefreshTokenMiddleware(cfg *config.APIConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		data, err := cfg.DBQueries.GetRefreshToken(ctx, token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if data.ExpiredAt.Before(time.Now()) {
			http.Error(w, "Refresh token expired", http.StatusUnauthorized)
			return
		}

		if data.RevokedAt.Valid {
			http.Error(w, "Refresh token revoked", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, RefreshTokenKey, data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetPolkaAPIKey(cfg *config.APIConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetPolkaToken(r.Header)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		if token != cfg.POLKA_KEY {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
