package middleware

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "request does not contain an access token"})
		c.Abort()
		return
	}
	claim, err := ValidateToken(tokenString)
	if err != nil {
		c.JSON(401, gin.H{"error": "here1" + err.Error()})
		c.Abort()
		return
	}
	c.Set("Id", claim.ID)
	c.Set("Email", claim.Email)
	c.Next()
}

type JWTClaim struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func ValidateToken(signedToken string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(
				"TODO: REPLACE THIS STRING WITH LEGIT SECRET CODE",
			), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return nil, err
	}
	fmt.Println("claims.Email", claims.Email)
	fmt.Println("claims.ID", claims.ID)
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return nil, err
	}
	return claims, err
}
