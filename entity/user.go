package entity

import (
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Username string    `json:"username,omitempty"`
	Email    string    `json:"email,omitempty"`
	Password string    `json:"password,omitempty"`
}

type ProfileInfoResp struct {
	Name               string `json:"name,omitempty"`
	Username           string `json:"username,omitempty"`
	Email              string `json:"email,omitempty"`
	CollectionsCreated uint32 `json:"collectionsCreated,omitempty"`
	CardsCreated       uint32 `json:"cardsCreated,omitempty"`
	CardsMastered      uint32 `json:"cardsMastered,omitempty"`
	CardsReviewing     uint32 `json:"cardsReviewing,omitempty"`
	CardsLearning      uint32 `json:"cardsLearning,omitempty"`
	// CollectionsLiked
	// CollectionsDisliked
	// CollectionsStarred
}

type UserCardStatistics struct {
	CardsCreated   uint32 `json:"cardsCreated,omitempty"`
	CardsMastered  uint32 `json:"cardsMastered,omitempty"`
	CardsReviewing uint32 `json:"cardsReviewing,omitempty"`
	CardsLearning  uint32 `json:"cardsLearning,omitempty"`
}
type UserCollectionStatistics struct {
	CollectionsCreated uint32 `json:"collectionsCreated,omitempty"`
}

type UserWithAuthToken struct {
	User  *User
	Token string
}

type UserRegistration struct {
	Username string    `json:"username,omitempty"`
	Email    string    `json:"email,omitempty"`
	Password string    `json:"password,omitempty"`
	Token    uuid.UUID `json:"token,omitempty"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserRegistration) HashEncryptPassword() error {
	encryptedPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(encryptedPwd)
	return nil
}

func (u *User) ComparePasswords(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRegistration) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashEncryptPassword(); err != nil {
		return err
	}
	// if u.PhoneNumber != nil {
	// 	*u.PhoneNumber = strings.TrimSpace(*u.PhoneNumber)
	// }
	// if u.Role != nil {
	// 	*u.Role = strings.ToLower(strings.TrimSpace(*u.Role))
	// }
	return nil
}
