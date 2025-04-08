package app

// PagedResult 分页查询结果
type PagedResult struct {
	List      interface{} `json:"list"`       // 数据列表
	Total     int         `json:"total"`      // 总记录数
	Page      int         `json:"page"`       // 当前页码
	PageSize  int         `json:"page_size"`  // 每页记录数
	TotalPage int         `json:"total_page"` // 总页数
}

// NewPagedResult 创建分页结果
func NewPagedResult(list interface{}, total, page, pageSize int) PagedResult {
	return PagedResult{
		List:      list,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: CalculateTotalPages(total, pageSize),
	}
}

// CalculateTotalPages 计算总页数
func CalculateTotalPages(total, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return (total + pageSize - 1) / pageSize
}
