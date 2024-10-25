package model

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password", validate:"required,min=8,max=64"`
}

type CreateUserRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=8,max=64"`
	FirstName string   `json:"first_name" validate:"required,min=2,max=64"`
	LastName  string   `json:"last_name" validate:"required,min=2,max=64"`
	Username  string   `json:"username" validate:"required,min=2,max=64"`
	DOB       Date     `json:"dob" validate:"required"`
	Address   string   `json:"address" validate:"required,min=2,max=128"`
	Phone     string   `json:"phone" validate:"required,min=10,max=10"`
	Roles     []string `json:"roles,omitempty"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=64"`
	Description string `json:"description" validate:"required,min=2,max=128"`
	Icon        string `json:"icon" validate:"required,min=2,max=64"`
	Color       string `json:"color" validate:"required,min=2,max=64"`
}
