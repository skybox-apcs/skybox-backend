package repositories

import (
	"context"

	"skybox-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	database   *mongo.Database
	collection string
}

// NewUserRepository creates a new instance of the userRepository
func NewUserRepository(db *mongo.Database, collection string) *userRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

// CreateUser creates a new user
func (ur *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	collection := ur.database.Collection(ur.collection)

	_, err := collection.InsertOne(ctx, user)

	return err
}

// GetUserByEmail retrieves a user by email
func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := ur.database.Collection(ur.collection)

	user := &models.User{}
	err := collection.FindOne(ctx, map[string]string{"email": email}).Decode(user)

	return user, err
}

// GetUserByID retrieves a user by ID
func (ur *userRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	collection := ur.database.Collection(ur.collection)

	user := &models.User{}
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(user)

	return user, err
}
