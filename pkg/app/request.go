package app

// CreateUserRequest defines the structure for creating a user
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Phone    string `json:"phone" binding:"omitempty"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest defines the structure for updating a user
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty"`
	Phone    string `json:"phone" binding:"omitempty"`
}

// GetUserRequest defines the structure for getting a user by ID
type GetUserRequest struct {
	ID int `json:"id" binding:"required"`
}

// DeleteUserRequest defines the structure for deleting a user
type DeleteUserRequest struct {
	ID int `json:"id" binding:"required"`
}
