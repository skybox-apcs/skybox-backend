package repositories

import (
	"context"
	"fmt"

	"skybox-backend/internal/api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// userTokenRepository is the interface for user token repository
type userTokenRepository struct {
	database   *mongo.Database
	collection string
}

// NewUserTokenRepository creates a new instance of the userTokenRepository
func NewUserTokenRepository(db *mongo.Database, collection string) *userTokenRepository {
	return &userTokenRepository{
		database:   db,
		collection: collection,
	}
}

// CreateUserToken creates a new user token
func (utr *userTokenRepository) CreateUserToken(ctx context.Context, userToken *models.UserToken) error {
	collection := utr.database.Collection(utr.collection)

	_, err := collection.InsertOne(ctx, userToken)
	if err != nil {
		return err
	}

	return nil
}

// FindUserToken retrieves a user token by token
func (utr *userTokenRepository) FindUserToken(ctx context.Context, token string) (*models.UserToken, error) {
	collection := utr.database.Collection(utr.collection)

	userToken := &models.UserToken{}
	err := collection.FindOne(ctx, bson.M{"token": token}).Decode(userToken)
	if err != nil {
		return nil, err
	}

	return userToken, nil
}

// GetUserTokenByID retrieves a user token by ID
func (utr *userTokenRepository) GetUserTokenByID(ctx context.Context, id string) (*models.UserToken, error) {
	collection := utr.database.Collection(utr.collection)

	userToken := &models.UserToken{}
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(userToken)
	if err != nil {
		return nil, err
	}

	return userToken, nil
}

// GetUserTokenByUserID retrieves a list user token by user ID
func (utr *userTokenRepository) GetUserTokenByUserID(ctx context.Context, userID string) (*[]models.UserToken, error) {
	collection := utr.database.Collection(utr.collection)

	userIDHex, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// Find all user tokens for the given user ID
	cursor, err := collection.Find(ctx, bson.M{"user_id": userIDHex})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// For every user token in the cursor, decode it into a UserToken struct
	// and append it to the userTokens slice
	var userTokens []models.UserToken
	for cursor.Next(ctx) {
		var userToken models.UserToken
		if err := cursor.Decode(&userToken); err != nil {
			return nil, err
		}
		userTokens = append(userTokens, userToken)
	}

	// Check for any errors that occurred during iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &userTokens, nil
}

// DeleteUserToken deletes a user token by token
func (utr *userTokenRepository) DeleteUserToken(ctx context.Context, token string) error {
	collection := utr.database.Collection(utr.collection)

	// Delete the user token
	fmt.Println("Deleting token:", token)
	deleteResult, err := collection.DeleteOne(ctx, bson.M{"token": token})
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 0 {
		return fmt.Errorf("no token found with the provided token")
	}

	return nil
}

// DeleteUserTokensByUserID deletes all user tokens by user ID
func (utr *userTokenRepository) DeleteUserTokensByUserID(ctx context.Context, userID string) error {
	collection := utr.database.Collection(utr.collection)

	userIDHex, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	deleteResult, err := collection.DeleteMany(ctx, bson.M{"user_id": userIDHex})
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 0 {
		return fmt.Errorf("no tokens found for the provided user ID")
	}

	return nil
}
