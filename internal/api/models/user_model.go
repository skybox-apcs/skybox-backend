package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUsers = "users"
)

// User struct encapsulates the user model
type User struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username             string             `bson:"username" json:"username"`
	Email                string             `bson:"email" json:"email"`
	PasswordHash         string             `bson:"password_hash" json:"-"`
	LastLoginAt          time.Time          `bson:"last_login_at" json:"last_login_at"`
	LastPasswordChangeAt time.Time          `bson:"last_password_change_at" json:"last_password_change_at"`
	RootFolderID         primitive.ObjectID `bson:"root_folder_id" json:"root_folder_id"` // The root folder ID for the user
	CreatedAt            time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt            time.Time          `bson:"updated_at" json:"updated_at"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	UpdateUserLastLogin(ctx context.Context, id string) error
}
