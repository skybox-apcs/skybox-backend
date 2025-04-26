package models

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`   // Access token
	RefreshToken string `json:"refresh_token"`  // Refresh token
	ID           string `json:"id"`             // User ID
	Username     string `json:"username"`       // Username
	Email        string `json:"email"`          // User email
	RootFolderID string `json:"root_folder_id"` // Root folder ID
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`  // Access token
	RefreshToken string `json:"refresh_token"` // Refresh token
	ID           string `json:"id"`            // User ID
	Username     string `json:"username"`      // Username
	Email        string `json:"email"`         // User email
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutResponse struct {
}
