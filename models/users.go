package models

import "github.com/google/uuid"

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func NewUser(username string) *User {
	return &User{
		ID:       uuid.New().String(),
		Username: username,
	}
}
