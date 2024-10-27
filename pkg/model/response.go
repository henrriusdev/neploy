package model

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
