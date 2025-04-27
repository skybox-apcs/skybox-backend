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
func (fr *folderRepository) CreateFolder(ctx context.Context, folder *models.Folder) (*models.Folder, error) {
	collection := fr.database.Collection(fr.collection)
	userID := ctx.Value("x-user-id-hex").(primitive.ObjectID) // Get userID from context x-user-id-hex saved before

	// Get the folder ID from the parent folder if it exists
	if folder.ParentFolderID != primitive.NilObjectID {
		folder, err := fr.GetFolderByID(ctx, folder.ParentFolderID.Hex())
		if err != nil {
			return nil, fmt.Errorf("parent folder not found")
		}

		// TODO: Implement sharing functionality later
		if folder.OwnerID != userID {
			return nil, fmt.Errorf("user does not have permission to create a folder in this parent folder")
		}
	}

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
	userID := ctx.Value("x-user-id-hex").(primitive.ObjectID) // Get userID from context x-user-id-hex saved before

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

	return nil, fmt.Errorf("folder not found or deleted")
}

// GetFolderParentIDByFolderID retrieves the parent folder ID of folder ID
func (fr *folderRepository) GetFolderParentIDByFolderID(ctx context.Context, folderID string) (string, error) {
	collection := fr.database.Collection(fr.collection)

	folderIDHex, err := primitive.ObjectIDFromHex(folderID)
	if err != nil {
		return "", err
	}

	// Find the folder by ID and isDeleted := false
	var folder models.Folder
	err = collection.FindOne(ctx, bson.M{"_id": folderIDHex, "is_deleted": false}).Decode(&folder)
	if err != nil {
		return "", err
	}

	// Get the Parent Folder ID
	parentFolderID := folder.ParentFolderID.Hex()
	return parentFolderID, nil
}

// GetFolderContents retrieves the contents of a folder by ID
func (fr *folderRepository) GetFolderListInFolder(ctx context.Context, folderID string) ([]*models.Folder, error) {
	collection := fr.database.Collection(fr.collection)

	// Get the current folder and check if the user has permission to access the folder
	_, err := fr.GetFolderByID(ctx, folderID)
	if err != nil {
		return nil, err
	}

	// Check if folderID is a valid ObjectID
	folderIDHex, err := primitive.ObjectIDFromHex(folderID)
	if err != nil {
		return nil, err
	}

	// TODO: Implement sharing functionality later
	// Get all folder contents where parent_folder_id matches the folderID
	cursor, err := collection.Find(ctx, bson.M{"parent_folder_id": folderIDHex, "is_deleted": false})
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

// GetFolderResponseListInFolder retrieves the folder responses in a folder by ID
// For the ownerID, it will get the Username and Email from the token
// SELECT * FROM folders f
// JOIN users u ON f.owner_id = u.id
// WHERE f.parent_folder_id = folderID AND f.is_deleted = false
func (fr *folderRepository) GetFolderResponseListInFolder(ctx context.Context, folderID string) ([]*models.FolderResponse, error) {
	collection := fr.database.Collection(fr.collection)

	// Decode the results into a slice of FileResponse
	var folderResponse []*models.FolderResponse

	// Get the current folder and check if the user has permission to access the folder
	_, err := fr.GetFolderByID(ctx, folderID)
	if err != nil {
		return nil, err
	}

	// Check if folderID is a valid ObjectID
	folderIDHex, err := primitive.ObjectIDFromHex(folderID)
	if err != nil {
		return nil, err
	}

	// Define the aggregation pipeline
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"parent_folder_id": folderIDHex,
				"is_deleted":       false,
			},
		},
		{
			"$lookup": bson.M{
				"from":         models.CollectionUsers, // The users collection
				"localField":   "owner_id",             // The field in the files collection
				"foreignField": "_id",                  // The field in the users collection
				"as":           "owner_details",        // The field to store the joined data
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$owner_details",
				"preserveNullAndEmptyArrays": true, // Optional: Keep files without matching users
			},
		},
		{
			"$project": bson.M{
				"id":               "$_id",
				"name":             "$name",
				"owner_id":         "$owner_id",
				"owner_user_name":  "$owner_details.username",
				"owner_email":      "$owner_details.email",
				"parent_folder_id": "$parent_folder_id",
				"stats":            "$stats",
				"created_at":       "$created_at",
				"updated_at":       "$updated_at",
			},
		},
	}

	// TODO: Implement sharing functionality later
	// Get all folder contents where parent_folder_id matches the folderID
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each document into a slice of Folder
	if cursor.All(ctx, &folderResponse); err != nil {
		return nil, err
	}

	return folderResponse, nil
}

// GetFileListInFolder retrieves the files in a folder by ID
func (fr *folderRepository) GetFileListInFolder(ctx context.Context, folderID string) ([]*models.File, error) {
	collection := fr.database.Collection(models.CollectionFiles)

	// Get the current folder and check if the user has permission to access the folder
	_, err := fr.GetFolderByID(ctx, folderID)
	if err != nil {
		return nil, err
	}

	// Check if folderID is a valid ObjectID
	folderIDHex, err := primitive.ObjectIDFromHex(folderID)
	if err != nil {
		return nil, err
	}

	// TODO: Implement sharing functionality later
	// Get all files in the folder where parent_folder_id matches the folderID
	cursor, err := collection.Find(ctx, bson.M{"parent_folder_id": folderIDHex, "is_deleted": false})
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

// GetFileResponseListInFolder retrieves the file responses in a folder by ID
// For the ownerID, it will get the Username and Email from the token
// SELECT *
// FROM files f
// JOIN users u ON f.owner_id = u.id
// WHERE f.parent_folder_id = folderID AND f.is_deleted = false
func (fr *folderRepository) GetFileResponseListInFolder(ctx context.Context, folderID string) ([]*models.FileResponse, error) {
	folderCollection := fr.database.Collection(models.CollectionFiles)

	// Decode the results into a slice of FileResponse
	var fileResponses []*models.FileResponse

	// Get the current folder and check if the user has permission to access the folder
	_, err := fr.GetFolderByID(ctx, folderID)
	if err != nil {
		return nil, err
	}

	// Check if folderID is a valid ObjectID
	folderIDHex, err := primitive.ObjectIDFromHex(folderID)
	if err != nil {
		return nil, err
	}

	// Define the aggregation pipeline
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"parent_folder_id": folderIDHex,
				"is_deleted":       false,
			},
		},
		{
			"$lookup": bson.M{
				"from":         models.CollectionUsers, // The users collection
				"localField":   "owner_id",             // The field in the files collection
				"foreignField": "_id",                  // The field in the users collection
				"as":           "owner_details",        // The field to store the joined data
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$owner_details",
				"preserveNullAndEmptyArrays": true, // Optional: Keep files without matching users
			},
		},
		{
			"$project": bson.M{
				"id":               "$_id",
				"name":             "$file_name",
				"owner_id":         "$owner_id",
				"owner_username":   "$owner_details.username",
				"owner_email":      "$owner_details.email",
				"parent_folder_id": "$parent_folder_id",
				"size":             "$size",
				"mime_type":        "$mime_type",
				"created_at":       "$created_at",
				"updated_at":       "$updated_at",
			},
		},
	}

	cursor, err := folderCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each document into a slice of FileResponse
	if cursor.All(ctx, &fileResponses); err != nil {
		return nil, err
	}

	return fileResponses, nil
}

func (fr *folderRepository) DeleteFolder(ctx context.Context, id string) error {
	collection := fr.database.Collection(fr.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Check if the folder is not root
	folder, err := fr.GetFolderByID(ctx, id)
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

func (fr *folderRepository) RenameFolder(ctx context.Context, id string, newName string) error {
	collection := fr.database.Collection(fr.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Check if the folder is not root
	folder, err := fr.GetFolderByID(ctx, id)
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

func (fr *folderRepository) MoveFolder(ctx context.Context, id string, newParentID string) error {
	collection := fr.database.Collection(fr.collection)
	userID := ctx.Value("x-user-id-hex").(primitive.ObjectID)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	newParentIDHex, err := primitive.ObjectIDFromHex(newParentID)
	if err != nil {
		return err
	}

	// Check if the folder is not root
	folder, err := fr.GetFolderByID(ctx, id)
	if err != nil {
		return err
	}
	if folder.IsRoot {
		return fmt.Errorf("cannot move root folder")
	}
	// TODO: Implement sharing functionality later
	if folder.OwnerID != userID {
		return fmt.Errorf("user does not have permission to move this folder")
	}

	// Check if the new parent folder ID is valid
	parentFolder, err := fr.GetFolderByID(ctx, newParentID)
	if err != nil {
		return err
	}
	if parentFolder.OwnerID != userID {
		return fmt.Errorf("user does not have permission to move this folder to the new parent folder")
	}

	// Update the parent folder ID
	_, err = collection.UpdateOne(ctx, bson.M{"_id": idHex}, bson.M{
		"$set": bson.M{
			"parent_folder_id": newParentIDHex,
		},
	})

	return err
}
