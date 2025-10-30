package model

type User struct {
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	CreatedAt string  `json:"createdAt"`
	DeletedAt *string `json:"deletedAt"`
}
