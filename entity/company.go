package entity

import (
	"github.com/google/uuid"
)

type CompanyPremiumStatus string

const (
	CompanyPremiumStatus_Active   CompanyPremiumStatus = "active"
	CompanyPremiumStatus_Inactive CompanyPremiumStatus = "inactive"
)

type Company struct {
	Id            uuid.UUID            `json:"id,omitempty"`
	Name          string               `json:"name,omitempty"`
	ReferralToken uuid.UUID            `json:"referral_token,omitempty"`
	PremiumStatus CompanyPremiumStatus `json:"premium_status,omitempty"`
}
