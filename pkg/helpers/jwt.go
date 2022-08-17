package helpers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/flash-cards-vocab/backend/entity"
	// "github.com/AleksK1NG/api-mc/config"
	// "github.com/AleksK1NG/api-mc/internal/models"
)

type Claims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.StandardClaims
}

func GenerateJWTToken(user *entity.User) (string, error) {
	// Register the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Email: user.Email,
		ID:    user.Id.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Register the JWT string
	tokenString, err := token.SignedString([]byte("TODO: REPLACE THIS STRING WITH LEGIT SECRET CODE"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
