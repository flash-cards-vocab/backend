package entity

import "github.com/google/uuid"

type AuthContext struct {
	UserId uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
}
