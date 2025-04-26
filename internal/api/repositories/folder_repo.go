package repositories

import (
	"context"
	"time"

	"skybox-backend/internal/api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type folderRepository struct {
	database   *mongo.Database
	collection string
}

// NewFolderRepository creates a new instance of the folderRepository
func NewFolderRepository(db *mongo.Database, collection string) *folderRepository {
	return &folderRepository{
		database:   db,
		collection: collection,
	}
}

// CreateFolder creates a new folder
func (fr *folderRepository) CreateFolder(ctx context.Context, folder *models.Folder) (*models.Folder, error) {
	collection := fr.database.Collection(fr.collection)

	// Create folder in database
	result, err := collection.InsertOne(ctx, folder)
	if err != nil {
		return nil, err
	}

	// Assign the ID to the folder object
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		folder.ID = oid
	}

	return folder, nil
}

// GetFolderByID retrieves a folder by ID
func (fr *folderRepository) GetFolderByID(ctx context.Context, id string) (*models.Folder, error) {
	collection := fr.database.Collection(fr.collection)

	folder := &models.Folder{}
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Find the folder by ID and isDeleted := false
	err = collection.FindOne(ctx, bson.M{"_id": idHex, "is_deleted": false}).Decode(folder)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

// GetFolderParentIDByFolderID retrieves the parent folder ID of folder ID
func (fr *folderRepository) GetFolderParentIDByFolderID(ctx context.Context, folderID string) (string, error) {
	collection := fr.database.Collection(fr.collection)

	folderIDHex, err := primitive.ObjectIDFromHex(folderID)
	if err != nil {
		return "", err
	}

	// Find the folder by ID and isDeleted := false
	err = collection.FindOne(ctx, bson.M{"_id": folderIDHex, "is_deleted": false}).Decode(folderIDHex)
	if err != nil {
		return "", err
	}

	// Get the Parent Folder ID
	parentFolderID := folderIDHex.Hex()
	return parentFolderID, nil
}

// GetFolderContents retrieves the contents of a folder by ID
func (fr *folderRepository) GetFolderContents(ctx context.Context, folderID string) ([]*any, error) {
	collection := fr.database.Collection(fr.collection)

	// Check if folderID is a valid ObjectID
	var folderIDHex = primitive.NilObjectID
	var err error
	if folderID != "" {
		folderIDHex, err = primitive.ObjectIDFromHex(folderID)
		if err != nil {
			return nil, err
		}
	}

	// Get all folder contents where parent_folder_id matches the folderID
	cursor, err := collection.Find(ctx, bson.M{"parent_folder_id": folderIDHex, "is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// TODO)): Get file too

	// Iterate through the cursor and decode each document into a slice of any type
	var contents []*any
	for cursor.Next(ctx) {
		var content any
		if err := cursor.Decode(&content); err != nil {
			return nil, err
		}
		contents = append(contents, &content)
	}

	return contents, nil
}

func (fr *folderRepository) DeleteFolder(ctx context.Context, id string) error {
	collection := fr.database.Collection(fr.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Soft delete the folder by setting IsDeleted to true and updating DeletedAt timestamp
	_, err = collection.UpdateOne(ctx, bson.M{"_id": idHex}, bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		},
	})

	return err
}
