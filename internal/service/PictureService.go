package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"botanical-api2/internal/models"
	"botanical-api2/internal/repository"
	"botanical-api2/pkg/setting"
	"botanical-api2/pkg/utils"
)

// PictureService 图片服务接口
type PictureService interface {
	// 上传用户头像
	UploadAvatar(userID int, file *multipart.FileHeader) (*models.UserAvatar, error)

	// 获取用户头像
	GetUserAvatar(userID int) (*models.UserAvatar, error)

	// 上传用户图片
	UploadUserPicture(userID int, file *multipart.FileHeader) (*models.UserPicture, error)

	// 获取用户所有图片
	GetUserPictures(userID int) ([]models.UserPicture, error)

	// 获取单张图片信息
	GetPictureByID(id int) (*models.UserPicture, error)

	// 删除用户图片
	DeletePicture(userID int, pictureID int) error
}

// PictureServiceImpl 图片服务实现
type PictureServiceImpl struct {
	Repositories *repository.Repositories
}

// NewPictureService 创建图片服务
func NewPictureService(repositories *repository.Repositories) PictureService {
	return &PictureServiceImpl{
		Repositories: repositories,
	}
}

// UploadAvatar 上传用户头像
func (s *PictureServiceImpl) UploadAvatar(userID int, file *multipart.FileHeader) (*models.UserAvatar, error) {
	// 检查文件类型
	if !utils.IsImageFile(file.Filename) {
		return nil, errors.New("只允许上传图片文件")
	}

	// 生成唯一文件名
	fileName := fmt.Sprintf("avatar_%d_%d%s", userID, time.Now().UnixNano(), filepath.Ext(file.Filename))
	filePath := filepath.Join(setting.UploadPath, "avatars", fileName)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	// 保存文件
	if err := utils.SaveUploadedFile(file, filePath); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 构建URL
	url := fmt.Sprintf("%s/uploads/avatars/%s", setting.ServerDomain, fileName)

	// 检查用户是否已有头像
	existingAvatar, err := s.Repositories.Avatar.GetByUserID(userID)
	if err == nil && existingAvatar != nil {
		// 更新现有头像信息
		existingAvatar.URL = url
		if err := s.Repositories.Avatar.Update(existingAvatar); err != nil {
			return nil, fmt.Errorf("更新头像信息失败: %w", err)
		}
		return existingAvatar, nil
	}

	// 创建新头像记录
	avatar := &models.UserAvatar{
		URL:    url,
		UserID: userID,
	}

	if err := s.Repositories.Avatar.Create(avatar); err != nil {
		return nil, fmt.Errorf("创建头像记录失败: %w", err)
	}

	return avatar, nil
}

// GetUserAvatar 获取用户头像
func (s *PictureServiceImpl) GetUserAvatar(userID int) (*models.UserAvatar, error) {
	return s.Repositories.Avatar.GetByUserID(userID)
}

// UploadUserPicture 上传用户图片
func (s *PictureServiceImpl) UploadUserPicture(userID int, file *multipart.FileHeader) (*models.UserPicture, error) {
	// 检查文件类型
	if !utils.IsImageFile(file.Filename) {
		return nil, errors.New("只允许上传图片文件")
	}

	// 生成唯一文件名
	fileName := fmt.Sprintf("picture_%d_%d%s", userID, time.Now().UnixNano(), filepath.Ext(file.Filename))
	filePath := filepath.Join(setting.UploadPath, "pictures", fileName)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	// 保存文件
	if err := utils.SaveUploadedFile(file, filePath); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 构建URL
	url := fmt.Sprintf("%s/uploads/pictures/%s", setting.ServerDomain, fileName)

	// 创建图片记录
	picture := &models.UserPicture{
		URL:    url,
		UserID: userID,
	}

	if err := s.Repositories.Picture.Create(picture); err != nil {
		return nil, fmt.Errorf("创建图片记录失败: %w", err)
	}

	return picture, nil
}

// GetUserPictures 获取用户所有图片
func (s *PictureServiceImpl) GetUserPictures(userID int) ([]models.UserPicture, error) {
	return s.Repositories.Picture.GetByUserID(userID)
}

// GetPictureByID 获取单张图片信息
func (s *PictureServiceImpl) GetPictureByID(id int) (*models.UserPicture, error) {
	return s.Repositories.Picture.GetByID(id)
}

// DeletePicture 删除用户图片
func (s *PictureServiceImpl) DeletePicture(userID int, pictureID int) error {
	// 获取图片信息
	picture, err := s.Repositories.Picture.GetByID(pictureID)
	if err != nil {
		return fmt.Errorf("获取图片信息失败: %w", err)
	}

	if picture == nil {
		return errors.New("图片不存在")
	}

	// 检查图片是否属于该用户
	if picture.UserID != userID {
		return errors.New("无权删除该图片")
	}

	// 删除文件（可选，也可以保留文件）
	// 从URL中提取文件名
	fileName := filepath.Base(picture.URL)
	filePath := filepath.Join(setting.UploadPath, "pictures", fileName)

	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	// 删除数据库记录
	return s.Repositories.Picture.Delete(pictureID)
}
