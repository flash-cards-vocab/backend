package entity

import (
	"github.com/google/uuid"
)

type UserCompanySubscriptionStatus string

const (
	UserCompanySubscriptionStatus_Active   UserCompanySubscriptionStatus = "active"
	UserCompanySubscriptionStatus_Inactive UserCompanySubscriptionStatus = "inactive"
)

type UserCompanySubscription struct {
	Id        uuid.UUID                     `json:"id,omitempty"`
	CompanyId uuid.UUID                     `json:"companyId,omitempty"`
	UserId    uuid.UUID                     `json:"userId,omitempty"`
	Status    UserCompanySubscriptionStatus `json:"status,omitempty"`
}
