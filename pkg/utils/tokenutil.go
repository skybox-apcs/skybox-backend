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
		"ID":       user.ID,
		"Email":    user.Email,
		"Username": user.Username,
		"exp":      exp,
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

// GetKeyFromToken gets the [Key] value from the token
func GetKeyFromToken(key string, requestToken string, secret string) (string, error) {
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

	if claims[key] == nil {
		return "", fmt.Errorf("key %s not found in token", key)
	}

	return claims[key].(string), nil
}

// GenerateToken generates a custom interface token
func GenerateToken(data map[string]string, secret string, expiry int) (string, error) {
	// Calculate the expiry date
	exp := time.Now().Add(time.Duration(expiry) * time.Hour).Unix()

	// Create the claims
	claims := jwt.MapClaims{
		"exp": exp,
	}

	for key, value := range data {
		claims[key] = value
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	customToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error while signing token: %w", err)
	}

	return customToken, nil
}

// GetKeysFromToken gets all claims from the token
func GetKeysFromToken(requestToken string, secret string) (map[string]string, error) {
	// Parse the token
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error while parsing token: %w", err)
	}

	// Get the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, fmt.Errorf("error while getting claims from token")
	}

	data := make(map[string]string)
	for key, value := range claims {
		switch v := value.(type) {
		case string:
			data[key] = v
		case float64:
			data[key] = fmt.Sprintf("%f", v)
		default:
			return nil, fmt.Errorf("unsupported type for key %s: %T", key, v)
		}
	}

	return data, nil
}
