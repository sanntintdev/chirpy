package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateChirpRequest struct {
	Body string `json:"body"`
}

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChirpValidResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserChirpyRedRequest struct {
	ChirpyRed string `json:"chirpy_red"`
}

type PolkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UpdateUserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
