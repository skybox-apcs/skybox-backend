package models

type NonResumableUploadResponse struct {
	FileID     string `json:"file_id"`     // The ID of the uploaded file
	FileName   string `json:"file_name"`   // The name of the uploaded file
	FileSize   int64  `json:"file_size"`   // The size of the uploaded file in bytes
	ChunkCount int    `json:"chunk_count"` // The number of chunks the file was split into (if applicable)
}

type AddChunkSessionRequest struct {
	ChunkNumber int    `json:"chunk_number"`
	ChunkSize   int    `json:"chunk_size"`
	ChunkHash   string `json:"chunk_hash"`
}
