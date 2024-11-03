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
