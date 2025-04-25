package utils

import (
	"fmt"
	"time"

	"skybox-backend/internal/api/models"

	"github.com/golang-jwt/jwt/v5"
)

// CreateAccessToken creates an access token for the user
func CreateAccessToken(user *models.User, secret string, expiry int) (string, error) {
	// Calculate the expiry date
	exp := time.Now().Add(time.Duration(expiry) * time.Hour).Unix()

	// Create the claims
	claims := jwt.MapClaims{
		"ID":    user.ID,
		"Email": user.Email,
		"exp":   exp,
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error while signing token: %w", err)
	}

	return accessToken, nil
}

// CreateRefreshToken creates a refresh token for the user
func CreateRefreshToken(user *models.User, secret string, expiry int) (string, error) {
	// Calculate the expiry date
	exp := time.Now().Add(time.Duration(expiry) * time.Hour).Unix()

	// Create the claims
	claims := jwt.MapClaims{
		"ID":    user.ID,
		"Email": user.Email,
		"exp":   exp,
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	refreshToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error while signing token: %w", err)
	}

	return refreshToken, nil
}

// IsAuthorized checks if the request is authorized
func IsAuthorized(requestToken string, secret string) (bool, error) {
	// Parse the token
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return false, fmt.Errorf("error while parsing token: %w", err)
	}

	// Check if the token is valid
	return token.Valid, nil
}

// GetIDFromToken gets the ID from the token
func GetIDFromToken(requestToken string, secret string) (string, error) {
	// Parse the token
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("error while parsing token: %w", err)
	}

	// Get the ID from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return "", fmt.Errorf("error while getting claims from token")
	}

	return claims["ID"].(string), nil
}
