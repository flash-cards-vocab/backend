package user_repository

import (
	"time"

	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID  `gorm:"primary_key;column:id"`
	Name      string     `gorm:"column:name"`
	Username  string     `gorm:"column:username"`
	Email     string     `gorm:"column:email"`
	Password  string     `gorm:"column:password"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (u *User) ToEntity() *entity.User {
	return &entity.User{
		Id:       u.Id,
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
}
