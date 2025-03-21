package service

import (
	"botanical-api2/internal/models"
	"botanical-api2/internal/repository"
	"botanical-api2/pkg/jwt"
	"botanical-api2/pkg/setting"
	"errors"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// RegisterUser方法，加密密码
func (s *UserService) RegisterUser(user *models.User) error {
	// 检查用户手机号是否存在
	exists, err := s.repo.CheckPhoneExists(user.Phone)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("用户已存在")
	}

	// 验证密码强度
	if err := validatePassword(user.PasswordHash); err != nil {
		return err
	}

	// 对密码进行哈希处理
	hashedPassword, err := generatePasswordHash(user.PasswordHash)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword
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

// UpdateUser 更新用户信息 - 使用动态字段更新
func (s *UserService) UpdateUser(user *models.User, fieldsToUpdate map[string]interface{}) error {
	// 检查用户是否存在
	_, err := s.repo.GetUserByID(user.ID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 调用仓库层的动态更新方法
	return s.repo.UpdateUserFields(user.ID, fieldsToUpdate)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}

// LoginUser 用户登录并生成令牌
func (s *UserService) LoginUser(phone, password string) (string, *models.User, error) {
	// 根据手机号获取用户
	user, err := s.repo.GetUserByPhone(phone)
	if err != nil {
		return "", nil, err // 修复: null 改为 nil
	}

	// 验证密码
	if !comparePasswords(user.PasswordHash, password) {
		return "", nil, errors.New("密码错误") // 修复: 添加缺失的用户返回值
	}

	// 生成JWT令牌
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", nil, err // 修复: 添加缺失的用户返回值
	}

	// 设置令牌相关信息
	expireTime := time.Now().Add(time.Duration(setting.JwtExpireHours) * time.Hour)

	// 创建更新字段映射
	updateFields := map[string]interface{}{
		"token":             token,
		"token_expire_time": &expireTime,
	}

	// 使用动态字段更新方法更新用户信息
	if err := s.repo.UpdateUserFields(user.ID, updateFields); err != nil {
		return "", nil, err // 修复: 添加缺失的用户返回值
	}

	// 更新本地user对象以反映数据库更改
	user.Token = token
	user.TokenExpireTime = &expireTime

	// 返回令牌和用户信息
	return token, user, nil // 修复: 返回完整的用户对象
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

// ChangePassword 修改用户密码
func (s *UserService) ChangePassword(userId int, oldPassword, newPassword string) error {

	// 验证密码强度
	if err := validatePassword(newPassword); err != nil {
		return err
	}

	// 获取用户信息
	user, err := s.repo.GetUserByID(userId)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !comparePasswords(user.PasswordHash, oldPassword) {
		return errors.New("原密码不正确")
	}

	// 对新密码进行哈希处理
	hashedPassword, err := generatePasswordHash(newPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	updateFields := map[string]interface{}{
		"password_hash": hashedPassword,
		// 不需要手动更新UpdatedAt，GORM会自动处理
	}

	return s.repo.UpdateUserFields(userId, updateFields)
}

// validatePassword 验证密码强度
func validatePassword(password string) error {
	if len(password) < 6 || len(password) > 20 {
		return errors.New("密码长度必须在6-20个字符之间")
	}

	// 可以添加更多密码强度验证
	hasNumber := false
	hasLetter := false

	for _, char := range password {
		if unicode.IsDigit(char) {
			hasNumber = true
		} else if unicode.IsLetter(char) {
			hasLetter = true
		}
	}

	if !hasNumber || !hasLetter {
		return errors.New("密码必须包含数字和字母")
	}

	return nil
}
