package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sanntintdev/chirpy/internal/auth"
	"github.com/sanntintdev/chirpy/internal/database"
	"github.com/sanntintdev/chirpy/pkg/config"
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
		}

		if err := auth.ComparePassword(user.HashedPassword, request.Password); err != nil {
			utils.RespondWithErr(w, http.StatusUnauthorized, "Invalid credentials", err)
		}

		utils.RespondWithJSON(w, http.StatusOK, &models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
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

		utils.RespondWithJSON(w, http.StatusCreated, models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
}
