package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sanntintdev/chirpy/internal/database"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secretKey []byte
	issuer    string
}

func NewTokenService(secretKey []byte, issuer string) *TokenService {
	return &TokenService{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

func (ts *TokenService) GenerateToken(user *database.User, expiration time.Duration) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ts.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.secretKey)
}

func (ts *TokenService) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return ts.secretKey, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return uuid.Nil, fmt.Errorf("invalid token")
}
