package models

type AddChunkRequest struct {
	ChunkNumber int `json:"chunk_number"` // The number of the chunk being uploaded
}

type AddChunkResponse struct {
}

type AddChunkViaFileIDRequest struct {
	ChunkNumber int `json:"chunk_number"` // The number of the chunk being uploaded
}

type AddChunkViaFileIDResponse struct {
}
