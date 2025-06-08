package repositories

import (
	"context"
	"time"

	"skybox-backend/internal/api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	database   *mongo.Database
	collection string
}

// NewUserRepository creates a new instance of the UserRepository
func NewUserRepository(db *mongo.Database, collection string) *UserRepository {
	return &UserRepository{
		database:   db,
		collection: collection,
	}
}

// CreateUser creates a new user
func (ur *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
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
func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := ur.database.Collection(ur.collection)

	user := &models.User{}
	err := collection.FindOne(ctx, map[string]string{"email": email}).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUsersByEmails retrieves users by their emails
func (ur *UserRepository) GetUsersByEmails(ctx context.Context, emails []string) ([]*models.User, error) {
	collection := ur.database.Collection(ur.collection)

	// Query users with emails in the list
	cursor, err := collection.Find(ctx, bson.M{"email": bson.M{"$in": emails}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode the results into a slice of users
	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID retrieves a user by ID
func (ur *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	collection := ur.database.Collection(ur.collection)

	user := &models.User{}
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(user)

	return user, err
}

// GetUsersByIDs retrieves users by their IDs
func (ur *UserRepository) GetUsersByIDs(ctx context.Context, ids []string) ([]*models.User, error) {
	collection := ur.database.Collection(ur.collection)

	// Convert string IDs to ObjectIDs
	objectIDs := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs[i] = objectID
	}

	// Query users with IDs in the list
	cursor, err := collection.Find(ctx, bson.M{"_id": bson.M{"$in": objectIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode the results into a slice of users
	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByUsername retrieves a user by username
func (ur *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	collection := ur.database.Collection(ur.collection)

	user := &models.User{}
	err := collection.FindOne(ctx, map[string]string{"username": username}).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserLastLogin updates the last login time of a user
func (ur *UserRepository) UpdateUserLastLogin(ctx context.Context, id string) error {
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
