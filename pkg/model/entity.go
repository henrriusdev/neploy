package model

type BaseEntity struct {
	ID        string `json:"id" db:"id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
	DeletedAt string `json:"deleted_at" db:"deleted_at"`
}

type User struct {
	BaseEntity
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	Email     string `json:"email" db:"email"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	DOB       string `json:"dob" db:"dob"`
	Address   string `json:"address" db:"address"`
	Phone     string `json:"phone" db:"phone"`
}

type Role struct {
	BaseEntity
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type Application struct {
	BaseEntity
}

type ApplicationUser struct {
	BaseEntity
}
