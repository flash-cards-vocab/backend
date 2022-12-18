package company_repository

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

func New(db *gorm.DB) repositoryIntf.CompanyRepository {
	return &repository{db: db, tableName: "users"}
}

func (r *repository) CreateUserCompanySubscription(userId, referralToken uuid.UUID) error {
	// Search for a company which referral token == referralToken
	var company *entity.Company
	err := r.db.Table("company").Where("referral_token=?", referralToken).First(&company).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	var subscriptionStatus entity.UserCompanySubscriptionStatus
	if company.PremiumStatus == entity.CompanyPremiumStatus_Inactive {
		subscriptionStatus = entity.UserCompanySubscriptionStatus_Inactive
	} else {
		subscriptionStatus = entity.UserCompanySubscriptionStatus_Active
	}

	subscription := &entity.UserCompanySubscription{
		Id:        uuid.New(),
		CompanyId: company.Id,
		UserId:    userId,
		Status:    subscriptionStatus,
	}
	err = r.db.Table("user_company_subscription").Create(&subscription).Error
	if err != nil {
		return err
	}
	return nil
}
