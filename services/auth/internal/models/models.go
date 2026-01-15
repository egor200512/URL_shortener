package models

import "github.com/google/uuid"

type RegisterRequest struct {
	Email    string
	Password string
}

type User struct {
	ID       uuid.UUID
	Email    string
	Password string
}
