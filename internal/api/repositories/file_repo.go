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

type fileRepository struct {
	database   *mongo.Database
	collection string
}

// NewFileRepository creates a new instance of the fileRepository
func NewFileRepository(db *mongo.Database, collection string) *fileRepository {
	return &fileRepository{
		database:   db,
		collection: collection,
	}
}

func (fr *fileRepository) UploadFileMetadata(ctx context.Context, file *models.File) (*models.File, error) {
	collection := fr.database.Collection(fr.collection)
	folderCollection := fr.database.Collection(models.CollectionFolders)
	userID := ctx.Value("x-user-id-hex").(primitive.ObjectID)

	// Get the folder
	var folder models.Folder
	err := folderCollection.FindOne(ctx, bson.M{"_id": file.ParentFolderID}).Decode(&folder)
	if err != nil {
		return nil, err
	}

	// Check if the folder belongs to the user
	if folder.OwnerID != userID {
		return nil, mongo.ErrNoDocuments
	}

	// Insert the file metadata into the database
	result, err := collection.InsertOne(ctx, file)
	if err != nil {
		return nil, err
	}

	// Set the ID of the file to the inserted ID
	file.ID = result.InsertedID.(primitive.ObjectID)

	return file, nil
}

func (fr *fileRepository) GetFileByID(ctx context.Context, id string) (*models.File, error) {
	collection := fr.database.Collection(fr.collection)
	userID := ctx.Value("x-user-id-hex").(primitive.ObjectID)

	file := &models.File{}
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid file ID: %v", err)
	}

	// Find the file by ID and isDeleted := false
	err = collection.FindOne(ctx, bson.M{"_id": idHex, "is_deleted": false, "owner_id": userID}).Decode(file)
	if err == nil {
		return file, nil
	}
	// TODO: Implement sharing functionality later
	// If the file is not found, check if it is shared with the user

	return nil, err
}

func (fr *fileRepository) DeleteFile(ctx context.Context, id string) error {
	collection := fr.database.Collection(fr.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Soft delete the file by setting is_deleted to true and deleted_at to current time
	_, err = collection.UpdateOne(ctx, bson.M{"_id": idHex}, bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

func (fr *fileRepository) RenameFile(ctx context.Context, id string, newName string) error {
	collection := fr.database.Collection(fr.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Check if the file related to the user (via owner or sharing)
	_, err = fr.GetFileByID(ctx, id)
	if err != nil {
		return err
	}

	// Rename the file by updating the file_name field
	_, err = collection.UpdateOne(ctx, bson.M{"_id": idHex}, bson.M{
		"$set": bson.M{
			"file_name": newName,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to rename file: %v", err)
	}

	return nil
}

func (fr *fileRepository) MoveFile(ctx context.Context, id string, newParentFolderID string) error {
	collection := fr.database.Collection(fr.collection)
	folderCollection := fr.database.Collection(models.CollectionFolders)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Check if the file related to the user (via owner or sharing)
	_, err = fr.GetFileByID(ctx, id)
	if err != nil {
		return err
	}

	newParentIDHex, err := primitive.ObjectIDFromHex(newParentFolderID)
	if err != nil {
		return fmt.Errorf("invalid new parent folder ID: %v", err)
	}

	// Check if the new parent folder ID is exist
	var folder models.Folder
	err = folderCollection.FindOne(ctx, bson.M{"_id": newParentIDHex}).Decode(&folder)
	if err != nil {
		return err
	}

	// Move the file by updating the parent_folder_id field
	_, err = collection.UpdateOne(ctx, bson.M{"_id": idHex}, bson.M{
		"$set": bson.M{
			"parent_folder_id": newParentIDHex,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to move file: %v", err)
	}

	return nil
}

func (fr *fileRepository) SearchFiles(ctx context.Context, ownerId primitive.ObjectID, query string) ([]*models.File, error) {
	collection := fr.database.Collection(fr.collection)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Define the filter for searching files
	filter := bson.M{
		"$and": []bson.M{
			{"owner_id": ownerId},
			{"file_name": bson.M{"$regex": query, "$options": "i"}}, // Case-insensitive regex match
			{"is_deleted": false}, // Exclude deleted files
			{"status": "uploaded"},
		},
	}

	// Execute the query
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to search files: %v", err)
	}
	defer cursor.Close(ctx)

	// Parse the results
	var files []*models.File
	for cursor.Next(ctx) {
		var file models.File
		if err := cursor.Decode(&file); err != nil {
			return nil, fmt.Errorf("failed to decode file: %v", err)
		}
		files = append(files, &file)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return files, nil
}
