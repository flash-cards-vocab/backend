package repository

import (
	"errors"

	"github.com/google/uuid"
)

var ErrCompanyNotFound = errors.New("Company not found")

type CompanyRepository interface {
	CreateUserCompanySubscription(userId, referralToken uuid.UUID) error
}
