// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID             uuid.UUID
	Email          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	HashedPassword string
}
