package repositories

import (
	"context"
	"fmt"
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
func (fr *folderRepository) CreateFolder(ctx context.Context, folder *models.Folder, userID primitive.ObjectID) (*models.Folder, error) {
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
func (fr *folderRepository) GetFolderByID(ctx context.Context, id string, userID primitive.ObjectID) (*models.Folder, error) {
	collection := fr.database.Collection(fr.collection)

	folder := &models.Folder{}
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Find the folder by ID and isDeleted := false
	// Owner priority
	err = collection.FindOne(ctx, bson.M{"_id": idHex, "is_deleted": false, "owner_id": userID}).Decode(folder)
	if err == nil {
		// If the folder is found, return it
		return folder, nil
	}

	// TODO: Implement sharing functionality later

	return folder, fmt.Errorf("folder not found or deleted")
}

// GetFolderParentIDByFolderID retrieves the parent folder ID of folder ID
func (fr *folderRepository) GetFolderParentIDByFolderID(ctx context.Context, folderID string, userID primitive.ObjectID) (string, error) {
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
func (fr *folderRepository) GetFolderListInFolder(ctx context.Context, folderID string, userID primitive.ObjectID) ([]*models.Folder, error) {
	collection := fr.database.Collection(fr.collection)

	// Check if folderID is a valid ObjectID
	folderIDHex, err := primitive.ObjectIDFromHex(folderID)
	if err != nil {
		return nil, err
	}

	// TODO: Implement sharing functionality later
	// Get all folder contents where parent_folder_id matches the folderID
	cursor, err := collection.Find(ctx, bson.M{"parent_folder_id": folderIDHex, "is_deleted": false, "owner_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each document into a slice of Folder
	var contents []*models.Folder
	for cursor.Next(ctx) {
		var folder models.Folder
		if err := cursor.Decode(&folder); err != nil {
			return nil, err
		}
		contents = append(contents, &folder)
	}

	return contents, nil
}

// GetFileListInFolder retrieves the files in a folder by ID
func (fr *folderRepository) GetFileListInFolder(ctx context.Context, folderID string, userId primitive.ObjectID) ([]*models.File, error) {
	collection := fr.database.Collection(models.CollectionFiles)

	// Check if folderID is a valid ObjectID
	folderIDHex, err := primitive.ObjectIDFromHex(folderID)
	if err != nil {
		return nil, err
	}

	// TODO: Implement sharing functionality later
	// Get all files in the folder where parent_folder_id matches the folderID
	cursor, err := collection.Find(ctx, bson.M{"parent_folder_id": folderIDHex, "is_deleted": false, "owner_id": userId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each document into a slice of File (or any other type)
	var files []*models.File
	for cursor.Next(ctx) {
		var file models.File // Replace with actual file type
		if err := cursor.Decode(&file); err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	return files, nil
}

func (fr *folderRepository) DeleteFolder(ctx context.Context, id string, userId primitive.ObjectID) error {
	collection := fr.database.Collection(fr.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Check if the folder is not root
	folder, err := fr.GetFolderByID(ctx, id, userId)
	if err != nil {
		return err
	}

	if folder.IsRoot {
		return fmt.Errorf("cannot delete root folder")
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

func (fr *folderRepository) RenameFolder(ctx context.Context, id string, newName string, userId primitive.ObjectID) error {
	collection := fr.database.Collection(fr.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Check if the folder is not root
	folder, err := fr.GetFolderByID(ctx, id, userId)
	if err != nil {
		return err
	}

	if folder.IsRoot {
		return fmt.Errorf("cannot rename root folder")
	}

	// Update the folder name
	_, err = collection.UpdateOne(ctx, bson.M{"_id": idHex}, bson.M{
		"$set": bson.M{
			"name": newName,
		},
	})

	return err
}

func (fr *folderRepository) MoveFolder(ctx context.Context, id string, newParentID string, userId primitive.ObjectID) error {
	collection := fr.database.Collection(fr.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	newParentIDHex, err := primitive.ObjectIDFromHex(newParentID)
	if err != nil {
		return err
	}

	// Check if the folder is not root
	folder, err := fr.GetFolderByID(ctx, id, userId)
	if err != nil {
		return err
	}
	if folder.IsRoot {
		return fmt.Errorf("cannot move root folder")
	}

	// Check if the new parent folder ID is valid
	_, err = fr.GetFolderByID(ctx, newParentID, userId)
	if err != nil {
		return err
	}

	// Update the parent folder ID
	_, err = collection.UpdateOne(ctx, bson.M{"_id": idHex}, bson.M{
		"$set": bson.M{
			"parent_folder_id": newParentIDHex,
		},
	})

	return err
}
