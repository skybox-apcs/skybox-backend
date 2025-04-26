package repositories

import (
	"context"
	"time"

	"skybox-backend/internal/api/models"

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
	folderCollection := ur.database.Collection(models.CollectionFolders)

	// Create the user in the database
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	// Create the root folder for the user
	rootFolder := &models.Folder{
		OwnerID:        result.InsertedID.(primitive.ObjectID),
		ParentFolderID: primitive.NilObjectID,
		Name:           "Root",
		IsDeleted:      false,
		IsRoot:         true,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	// Insert the root folder into the database
	folderResult, err := folderCollection.InsertOne(ctx, rootFolder)
	if err != nil {
		return err
	}

	// Update the user's root folder ID
	user.RootFolderID = folderResult.InsertedID.(primitive.ObjectID)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": result.InsertedID}, bson.M{"$set": bson.M{"root_folder_id": user.RootFolderID}})

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

// GetUserByUsername retrieves a user by username
func (ur *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	collection := ur.database.Collection(ur.collection)

	user := &models.User{}
	err := collection.FindOne(ctx, map[string]string{"username": username}).Decode(user)

	return user, err
}

// UpdateUserLastLogin updates the last login time of a user
func (ur *userRepository) UpdateUserLastLogin(ctx context.Context, id string) error {
	collection := ur.database.Collection(ur.collection)

	// Convert the string ID to ObjectID
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Update the last login time
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": idHex},
		bson.M{"$set": bson.M{"last_login_at": time.Now()}},
	).Err()
	if err != nil {
		return err
	}

	return nil
}
