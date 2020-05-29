package models

import (
	"gallery0api/rand"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserTable struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Name     string
	Token    string `gorm:"unique"`
}

type UserService interface {
	CreateUser(user *UserTable) error
	Login(user *UserTable) (string, error)
}

var _ UserService = &UserGorm{}

type UserGorm struct {
	db *gorm.DB
}

func NewUserGorm(db *gorm.DB) UserService {
	return &UserGorm{db}
}

func (ug *UserGorm) CreateUser(user *UserTable) error {
	return ug.db.Create(user).Error
}

var err error

func (ug *UserGorm) Login(user *UserTable) (string, error) {
	found := new(UserTable)
	if err := ug.db.Where("email = ?", user.Email).First(&found).Error; err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(user.Password))
	if err != nil {
		return "", nil
	}
	token, err := rand.GetToken()
	if err != nil {
		return "", err
	}
	err = ug.db.Model(&UserTable{}).
		Where("id = ?", found.ID).
		Update("token", token).Error
	if err != nil {
		return "", err
	}
	return token, nil
}
