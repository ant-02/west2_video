package repository

import (
	"west2/pkg/model"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByUsername(username string) (*model.User, error)
	GetUserById(id string) (*model.User, error)
	SetAvatar(id string, url string) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) CreateUser(user *model.User) error {
	return ur.db.Create(user).Error
}
func (ur *userRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := ur.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) GetUserById(id string) (*model.User, error) {
	var user model.User
	if err := ur.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) SetAvatar(id string, url string) error {
	return ur.db.Model(&model.User{}).Where("id = ?", id).Update("avatar", url).Error
}
