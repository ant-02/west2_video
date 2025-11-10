package service

import (
	"errors"
	"log"
	"west2/pkg/model"
	"west2/pkg/repository"
	"west2/util"

	"gorm.io/gorm"
)

type userService struct {
	ur repository.UserRepository
}

type UserService interface {
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{ur: ur}
}

func (us *userService) Login(username, password, code string) (*model.User, error) {
	// 根据用户名获取用户信息，判断用户是否存在
	user, err := us.ur.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("failed to get user from repository: username: %s, error: %v", username, err)
		return nil, err
	}

	// 进行密码比对
	if util.CheckPassword(password, user.Password) {
		user.Password = ""
		return user, nil
	}

	return nil, nil
}

func (us *userService) Register(username, password string) (bool, error) {
	// 判断用户名和密码是否为空
	if username == "" || password == "" {
		return false, errors.New("username or password is empty string")
	}

	// 判断用户是否已经被创建
	_, err := us.ur.GetUserByUsername(username)
	if err == nil {
		return false, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("failed to get user from repository: username: %s, error: %v", username, err)
		return false, err
	}

	// 对密码进行加密
	hash, err := util.HashPassword(password)
	if err != nil {
		log.Printf("failed to encode password: username: %s, error: %v", username, err)
		return false, err
	}

	// 雪花算法获取用户id
	id := util.GetID()

	// 将用户信息存入数据库
	user := &model.User{
		Id:       id,
		Username: username,
		Password: hash,
	}
	if err = us.ur.CreateUser(user); err != nil {
		log.Printf("failed to create a user: user: %+v, error: %v", user, err)
		return false, err
	}

	return true, nil
}

func (us *userService) GetUserInfoById(id string) (*model.User, error) {
	user, err := us.ur.GetUserById(id)
	if err != nil {
		log.Printf("failed to get user information by id: id: %s, error: %v", id, err)
		return nil, err
	}

	return user, nil
}

func (us *userService) UploadAvatar(id string, data string) (bool, error) {
	if err := util.SaveBase64Image(data, "./static/"+id+".png"); err != nil {
		log.Printf("failed to save image file: id: %s, error: %v", id, err)
		return false, err
	}

	if err := us.ur.SetAvatar(id, "http://localhost:8080/static/"+id+".png"); err != nil {
		log.Printf("failed to set user's avatar url: id: %s, error: %v", id, err)
		return false, err
	}

	return true, nil
}
