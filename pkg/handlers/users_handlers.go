package handlers

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/sanntintdev/chirpy/internal/auth"
	"github.com/sanntintdev/chirpy/internal/database"
	"github.com/sanntintdev/chirpy/pkg/config"
	"github.com/sanntintdev/chirpy/pkg/middleware"
	"github.com/sanntintdev/chirpy/pkg/models"
	"github.com/sanntintdev/chirpy/pkg/utils"
)

func LoginUserHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := &models.LoginUserRequest{}

		if err := utils.DecodeJSON(w, r, request); err != nil {
			return
		}

		ctx := r.Context()
		user, err := cfg.DBQueries.GetUserByEmail(ctx, request.Email)
		if err != nil {
			utils.RespondWithErr(w, http.StatusNotFound, "User not found", err)
			return
		}

		if err := auth.ComparePassword(user.HashedPassword, request.Password); err != nil {
			utils.RespondWithErr(w, http.StatusUnauthorized, "Invalid credentials", err)
			return
		}

		tokenService := auth.NewTokenService([]byte(cfg.JWT_SECRET_KEY), cfg.JWT_ISSUER)
		token, err := tokenService.GenerateToken(user)

		if err != nil {
			utils.RespondWithErr(w, http.StatusInternalServerError, "Failed to generate token", err)
			return
		}

		rfToken, err := auth.GenerateRefreshToken()
		if err != nil {
			utils.RespondWithErr(w, http.StatusInternalServerError, "Failed to generate refresh token", err)
			return
		}

		_, err = cfg.DBQueries.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
			UserID: user.ID,
			Token:  rfToken,
		})

		if err != nil {
			utils.RespondWithErr(w, http.StatusInternalServerError, "Failed to create refresh token", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, &models.UserResponse{
			ID:           user.ID,
			Email:        user.Email,
			Token:        token,
			RefreshToken: rfToken,
			IsChirpyRed:  user.IsChirpyRed,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
		})
	}
}

func CreateUserHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := &models.CreateUserRequest{}
		if err := utils.DecodeJSON(w, r, request); err != nil {
			return
		}

		hashedPassword, err := auth.HashPassword(request.Password)
		if err != nil {
			utils.RespondWithErr(w, http.StatusInternalServerError, "Failed to hash password", err)
			return
		}

		ctx := r.Context()
		user, err := cfg.DBQueries.CreateUser(ctx, database.CreateUserParams{
			ID:             uuid.New(),
			Email:          request.Email,
			HashedPassword: hashedPassword,
		})

		if err != nil {
			utils.RespondWithErr(w, http.StatusInternalServerError, "Failed to create user", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, models.UpdateUserResponse{
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
	}
}

func UpdateUserHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := ctx.Value(middleware.UserIDKey).(uuid.UUID)

		request := &models.UpdateUserRequest{}
		if err := utils.DecodeJSON(w, r, request); err != nil {
			return
		}

		user, err := cfg.DBQueries.GetUserByID(ctx, userID)
		if err != nil {
			utils.RespondWithErr(w, http.StatusNotFound, "User not found", err)
			return
		}

		if request.Email != "" {
			user.Email = request.Email
		}

		if request.Password != "" {
			hashedPassword, err := auth.HashPassword(request.Password)
			if err != nil {
				utils.RespondWithErr(w, http.StatusInternalServerError, "Failed to hash password", err)
				return
			}
			user.HashedPassword = hashedPassword
		}

		user, err = cfg.DBQueries.UpdateUser(ctx, database.UpdateUserParams{
			ID:             userID,
			Email:          user.Email,
			HashedPassword: user.HashedPassword,
		})

		if err != nil {
			utils.RespondWithErr(w, http.StatusInternalServerError, "Failed to update user", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, models.UpdateUserResponse{
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
	}
}

func PolkaWebhookHandler(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := &models.PolkaWebhookRequest{}
		if err := utils.DecodeJSON(w, r, request); err != nil {
			return
		}

		if request.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		ctx := r.Context()
		_, err := cfg.DBQueries.UpdateChirpyRedById(ctx, database.UpdateChirpyRedByIdParams{
			ID:         request.Data.UserID,
			IsChirpyRed: true,
		})

		if err != nil {
			if err == sql.ErrNoRows {
				utils.RespondWithErr(w, http.StatusNotFound, "User not found", err)
				return
			}
			utils.RespondWithErr(w, http.StatusInternalServerError, "Failed to update user", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
