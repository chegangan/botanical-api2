package util

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Pagination 分页结构体
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// NewPagination 创建新的分页对象
func NewPagination(c *gin.Context) *Pagination {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10 // 默认每页10条
	}

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// CalculateTotalPages 计算总页数
func CalculateTotalPages(totalItems int, pageSize int) int {
	if pageSize == 0 {
		return 0
	}
	return int(math.Ceil(float64(totalItems) / float64(pageSize)))
}

// SetPaginationHeader 设置分页响应头
func SetPaginationHeader(c *gin.Context, totalItems int, pagination *Pagination) {
	totalPages := CalculateTotalPages(totalItems, pagination.PageSize)
	c.Header("X-Total-Count", strconv.Itoa(totalItems))
	c.Header("X-Total-Pages", strconv.Itoa(totalPages))
	c.Header("X-Current-Page", strconv.Itoa(pagination.Page))
	c.Header("X-Page-Size", strconv.Itoa(pagination.PageSize))
}
