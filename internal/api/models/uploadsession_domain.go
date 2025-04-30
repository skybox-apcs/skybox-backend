package models

type AddChunkRequest struct {
	ChunkNumber int    `json:"chunk_number"` // The number of the chunk being uploaded
	ChunkSize   int    `json:"chunk_size"`   // The size of the chunk being uploaded
	ChunkHash   string `json:"chunk_hash"`   // The hash of the chunk being uploaded
}

type AddChunkResponse struct {
}

type AddChunkViaFileIDRequest struct {
	ChunkNumber int    `json:"chunk_number"` // The number of the chunk being uploaded
	ChunkSize   int    `json:"chunk_size"`   // The size of the chunk being uploaded
	ChunkHash   string `json:"chunk_hash"`   // The hash of the chunk being uploaded
}

type AddChunkViaFileIDResponse struct {
}
