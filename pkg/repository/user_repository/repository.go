package user_repository

import (
	"errors"

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
	err := r.db.Table(r.tableName).Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) CheckIfUserExistsByEmail(email string) (bool, error) {
	var user *entity.User
	err := r.db.Table(r.tableName).Where("email=?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}
	return true, nil
}

func (r *repository) GetUserByEmail(email string) (*entity.User, error) {
	var user *entity.User
	err := r.db.Table(r.tableName).Where("email=?", email).Find(&user).Error
	if err != nil {
		// if err == gorm.ErrRecordNotFound {
		// 	return nil, err
		// }
		return nil, err
	}
	return user, nil
}

func (r *repository) GetUserById(id uuid.UUID) (*entity.User, error) {
	var user *entity.User
	err := r.db.Table(r.tableName).Where("id=?", id).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
