package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUserTokens = "user_tokens"
)

// UserToken struct encapsulates the user token model
type UserToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`       // Reference to the user
	Token     string             `bson:"token" json:"token"`           // The token itself
	UserAgent string             `bson:"user_agent" json:"user_agent"` // User agent string
	IPAddress string             `bson:"ip_address" json:"ip_address"` // IP address of the user
	ExpiredAt time.Time          `bson:"expired_at" json:"expired_at"` // Expiration timestamp
	CreatedAt time.Time          `bson:"created_at" json:"created_at"` // Creation timestamp
}

type UserTokenRepository interface {
	CreateUserToken(ctx context.Context, userToken *UserToken) error               // Create a new user token
	FindUserToken(ctx context.Context, token string) (*UserToken, error)           // Retrieve a user token by token
	GetUserTokenByID(ctx context.Context, id string) (*UserToken, error)           // Retrieve a user token by ID
	GetUserTokenByUserID(ctx context.Context, userID string) (*[]UserToken, error) // Retrieve the list of tokens for a user
	DeleteUserToken(ctx context.Context, token string) error                       // Delete a user token by token
	DeleteUserTokensByUserID(ctx context.Context, userID string) error             // This is for log out every devices
}
