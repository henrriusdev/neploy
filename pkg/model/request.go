package model

import "time"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type CreateUserRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=8,max=64"`
	FirstName string   `json:"firstName" validate:"required,min=2,max=64"`
	LastName  string   `json:"lastName" validate:"required,min=2,max=64"`
	Username  string   `json:"username" validate:"required,min=2,max=64"`
	DOB       Date     `json:"dob" validate:"required"`
	Address   string   `json:"address" validate:"required,min=2,max=128"`
	Phone     string   `json:"phone" validate:"required,min=10,max=10"`
	Provider  string   `json:"provider" validate:"required,min=2,max=64"`
	Roles     []string `json:"roles,omitempty"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=64"`
	Description string `json:"description" validate:"required,min=2,max=128"`
	Icon        string `json:"icon" validate:"required,min=2,max=64"`
	Color       string `json:"color" validate:"required,min=2,max=64"`
}

type OnboardRequest struct {
	AdminUser CreateUserRequest   `json:"adminUser" validate:"required"`
	Roles     []CreateRoleRequest `json:"roles" validate:"required"`
	Metadata  MetadataRequest     `json:"metadata" validate:"required"`
	OauthID   string              `json:"-"`
}

type MetadataRequest struct {
	Name    string `json:"teamName" db:"team_name"`
	LogoURL string `json:"logoUrl" db:"logo_url"`
}

type InviteUserRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type CompleteInviteRequest struct {
	Token     string    `json:"token" validate:"required"`
	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName" validate:"required"`
	DOB       time.Time `json:"dob" validate:"required"`
	Phone     string    `json:"phone" validate:"required,min=10,max=10"`
	Address   string    `json:"address" validate:"required,min=2,max=128"`
	Email     string    `json:"email" validate:"required,email"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required,min=8,max=64"`
	OauthID   string    `json:"-"`
}

type CreateApplicationRequest struct {
	AppName     string `json:"appName"`
	Description string `json:"description"`
	TechStack   string `json:"techStack"`
}

type DeployApplicationRequest struct {
	RepoURL string `json:"repoUrl"`
}
