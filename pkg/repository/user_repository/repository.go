package user_repository

import (
	"errors"
	"time"

	repositoryIntf "github.com/flash-cards-vocab/backend/app/repository"
	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db        *gorm.DB
	tableName string
}

func New(db *gorm.DB) repositoryIntf.UserRepository {
	return &repository{db: db, tableName: "users"}
}

func (r *repository) CreateUser(user entity.User) (*entity.User, error) {
	user.Id = uuid.New()
	userToCreate := &User{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := r.db.Table(r.tableName).Create(&userToCreate).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) CheckIfUserExistsByEmail(email string) (bool, error) {
	var user *User
	err := r.db.Table(r.tableName).Where("email=? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}
	return true, nil
}

func (r *repository) CheckIfUsernameExists(username string) (bool, error) {
	var user *User
	err := r.db.Table(r.tableName).Where("username=? AND deleted_at IS NULL", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}
	return true, nil
}

func (r *repository) GetUserByEmail(email string) (*entity.User, error) {
	var user *User
	err := r.db.Table(r.tableName).Where("email=? AND deleted_at IS NULL", email).Find(&user).Error
	if err != nil {
		// if err == gorm.ErrRecordNotFound {
		// 	return nil, err
		// }
		return nil, err
	}
	return user.ToEntity(), nil
}

func (r *repository) GetUserById(id uuid.UUID) (*entity.User, error) {
	var user *User
	err := r.db.Table(r.tableName).Where("id=? AND deleted_at IS NULL", id).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user.ToEntity(), nil
}

func (r *repository) GetUserByUsername(username string) (*entity.User, error) {
	var user *User
	err := r.db.
		Table(r.tableName).
		Where("username=? AND deleted_at IS NULL", username).
		Find(&user).
		Error
	if err != nil {
		return nil, err
	}
	return user.ToEntity(), nil
}
