package models

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Count int64          `json:"count"`
}

type UserIDListRequest struct {
	IDs []string `json:"ids" binding:"required"`
}

type UserEmailListRequest struct {
	Emails []string `json:"emails" binding:"required"`
}
