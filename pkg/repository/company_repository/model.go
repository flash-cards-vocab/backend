package company_repository

import (
	"time"

	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
)

type Company struct {
	Id            uuid.UUID  `gorm:"primary_key;column:id"`
	Name          string     `gorm:"column:name"`
	ReferralToken uuid.UUID  `gorm:"column:referral_token"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
}

func (u *Company) ToEntity() *entity.Company {
	return &entity.Company{
		Id:            u.Id,
		Name:          u.Name,
		ReferralToken: u.ReferralToken,
	}
}

type UserCompanySubscription struct {
	Id        uuid.UUID                            `gorm:"primary_key;column:id"`
	CompanyId uuid.UUID                            `gorm:"column:company_id"`
	UserId    uuid.UUID                            `gorm:"column:user_id"`
	Status    entity.UserCompanySubscriptionStatus `gorm:"column:status"`
	CreatedAt time.Time                            `gorm:"column:created_at"`
	UpdatedAt time.Time                            `gorm:"column:updated_at"`
	DeletedAt *time.Time                           `gorm:"column:deleted_at"`
}

func (u *UserCompanySubscription) ToEntity() *entity.UserCompanySubscription {
	return &entity.UserCompanySubscription{
		Id:        u.Id,
		CompanyId: u.CompanyId,
		UserId:    u.UserId,
		Status:    u.Status,
	}
}
