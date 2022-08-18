package entity

import "github.com/google/uuid"

type AuthContext struct {
	UserId uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}
