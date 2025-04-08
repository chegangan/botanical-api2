package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// 允许的图片扩展名
var allowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".webp": true,
}

// IsImageFile 检查文件是否为图片
func IsImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return allowedImageExtensions[ext]
}

// SaveUploadedFile 保存上传的文件
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
