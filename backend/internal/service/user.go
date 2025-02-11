package service

import (
	"database/sql"
	"errors"
	"fmt"
	"mini-assets/internal/model"
	"mini-assets/internal/repository"
	"mini-assets/internal/util"
)

// UserService 定义了用户相关的业务逻辑接口
type UserService interface {
	RegisterUser(user *model.User) error
	LoginUser(username, password string) (string, error)
	GetUserByUsername(username string) (*model.User, error)
}

// userService 是UserService接口的实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建并返回一个UserService实例
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// RegisterUser 注册一个新的用户
func (s *userService) RegisterUser(user *model.User) error {
	// 检查用户名是否已存在
	existingUser, err := s.userRepo.GetByUsername(user.Username)
	if err != nil && err != sql.ErrNoRows {
		fmt.Printf("failed to check user existence: %v\n", err)
		return err
	}
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// 加密密码
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// 创建用户
	return s.userRepo.Create(user)
}

// LoginUser 用户登录
func (s *userService) LoginUser(username, password string) (string, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", err
	}

	// 验证密码
	if !util.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid password")
	}

	// 生成JWT令牌
	token, err := util.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByUsername 根据用户名获取用户信息
func (s *userService) GetUserByUsername(username string) (*model.User, error) {
	return s.userRepo.GetByUsername(username)
}
