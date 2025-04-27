package repositories

import (
	"context"
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

	file := &models.File{}
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Find the file by ID and isDeleted := false
	err = collection.FindOne(ctx, bson.M{"_id": idHex, "is_deleted": false}).Decode(file)
	if err != nil {
		return nil, err
	}

	return file, nil
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
		return err
	}

	return nil
}

func (fr *fileRepository) RenameFile(ctx context.Context, id string, newName string) error {
	collection := fr.database.Collection(fr.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
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
		return err
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

	newParentIDHex, err := primitive.ObjectIDFromHex(newParentFolderID)
	if err != nil {
		return err
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
		return err
	}

	return nil
}
