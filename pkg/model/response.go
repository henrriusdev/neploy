package model

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type OAuthResponse struct {
	Provider Provider `json:"provider"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
}

type UserResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

type TeamResponse struct {
	UserRoles     []UserRoles     `json:"userRoles"`
	UserTechStack []UserTechStack `json:"userTechStack"`
}
