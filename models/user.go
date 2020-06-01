package models

import (
	"encoding/base64"
	"fmt"
	"gallery0api/rand"
	"hash"

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
	GetByToken(token string) (*UserTable, error)
}

var _ UserService = &UserGorm{}

type UserGorm struct {
	db   *gorm.DB
	hmac hash.Hash
}

func NewUserGorm(db *gorm.DB, hmac hash.Hash) UserService {
	return &UserGorm{db, hmac}
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
	fmt.Println("token", token)
	ug.hmac.Write([]byte(token))
	hash := ug.hmac.Sum(nil)
	ug.hmac.Reset()
	fmt.Println("hash", base64.URLEncoding.EncodeToString(hash))
	encode := base64.URLEncoding.EncodeToString(hash)
	err = ug.db.Model(&UserTable{}).
		Where("id = ?", found.ID).
		Update("token", encode).Error
	if err != nil {
		return "", err
	}
	return token, nil
}

func (ug *UserGorm) GetByToken(token string) (*UserTable, error) {
	user := &UserTable{}
	err := ug.db.Where("token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
