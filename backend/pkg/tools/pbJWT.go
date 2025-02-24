package tools

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// PocketBaseJWTPayload represents the custom JWT claims
type PocketBaseJWTPayload struct {
	CollectionId string `json:"collectionId"`
	Id           string `json:"id"`
	Exp          int    `json:"exp"`
	Type         string `json:"type"`
	Refreshable  bool   `json:"refreshable"`
	jwt.RegisteredClaims
}

// GetUserIdFromJwt extracts the user ID from a JWT token without verification
func GetUserIdFromJwt(tokenString string) (string, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &PocketBaseJWTPayload{})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*PocketBaseJWTPayload); ok {
		return claims.Id, nil
	}

	return "", errors.New("invalid token or claims")
}

// VerifyJWTExpiration checks if the JWT token is still valid based on its expiration time.
func VerifyJWTExpiration(tokenString string) (bool, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &PocketBaseJWTPayload{})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(*PocketBaseJWTPayload); ok {
		return time.Now().Unix() < int64(claims.Exp), nil
	}
	return false, errors.New("invalid token or claims")
}
