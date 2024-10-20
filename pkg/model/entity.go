package model

type BaseEntity struct {
	ID        string `json:"id" db:"id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
	DeletedAt string `json:"deleted_at" db:"deleted_at"`
}

type User struct {
	BaseEntity
}

type Role struct {
	BaseEntity
}

type Application struct {
	BaseEntity
}

type ApplicationUser struct {
	BaseEntity
}
