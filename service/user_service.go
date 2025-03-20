package service

import (
	"botanical-api2/models"
	"botanical-api2/pkg/jwt"
	"botanical-api2/pkg/setting"
	"botanical-api2/repository"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		repo: repository.NewUserRepository(db),
	}
}

// 修改RegisterUser方法，加密密码
func (s *UserService) RegisterUser(user *models.User) error {
	// 检查用户是否存在
	exists, err := s.repo.CheckUserExist(user.Username, user.Phone)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("用户已存在")
	}

	// 对密码进行哈希处理
	hashedPassword, err := generatePasswordHash(user.PasswordHash)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return s.repo.CreateUser(user)
}

// CreateUser 创建用户（保持与路由处理函数的一致性）
func (s *UserService) CreateUser(user *models.User) error {
	// 调用已有的 RegisterUser 方法实现用户创建逻辑
	return s.RegisterUser(user)
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now()
	return s.repo.UpdateUser(user)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}

// 修改LoginUser方法，生成令牌
func (s *UserService) LoginUser(username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	// 验证密码
	if !comparePasswords(user.PasswordHash, password) {
		return "", errors.New("密码错误")
	}

	// 使用新的 jwt 包生成令牌
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	// 更新用户的令牌和过期时间
	user.Token = token
	user.TokenExpireTime = time.Now().Add(time.Duration(setting.JwtExpireHours) * time.Hour)
	err = s.repo.UpdateUser(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// generatePasswordHash 生成密码哈希
func generatePasswordHash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// comparePasswords 比较密码哈希
func comparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
