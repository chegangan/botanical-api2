package models

import "strings"

// ParseImageList 解析逗号分隔的图片URL列表
func ParseImageList(imageListStr string) []string {
	if imageListStr == "" {
		return []string{}
	}
	return strings.Split(imageListStr, ",")
}

// JoinImageList 将图片URL数组合并为逗号分隔的字符串
func JoinImageList(images []string) string {
	return strings.Join(images, ",")
}
