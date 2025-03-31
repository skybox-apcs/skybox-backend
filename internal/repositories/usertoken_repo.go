package repositories

import (
	"context"

	"skybox-backend/internal/models"

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

	cursor, err := collection.Find(ctx, bson.M{"user_id": userIDHex})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userTokens []models.UserToken
	for cursor.Next(ctx) {
		var userToken models.UserToken
		if err := cursor.Decode(&userToken); err != nil {
			return nil, err
		}
		userTokens = append(userTokens, userToken)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &userTokens, nil
}

// DeleteUserToken deletes a user token by ID
func (utr *userTokenRepository) DeleteUserToken(ctx context.Context, id string) error {
	collection := utr.database.Collection(utr.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": idHex})
	if err != nil {
		return err
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

	_, err = collection.DeleteMany(ctx, bson.M{"user_id": userIDHex})
	if err != nil {
		return err
	}

	return nil
}
