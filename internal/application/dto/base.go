package dto

// PageRequest 分页请求
type PageRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

// GetOffset 获取偏移量
func (p *PageRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.GetPageSize()
}

// GetPageSize 获取每页大小
func (p *PageRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return p.PageSize
}

// PageResponse 分页响应
type PageResponse struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(data interface{}) *Response {
	return &Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// SuccessWithPage 带分页的成功响应
func SuccessWithPage(data interface{}, page, pageSize int, total int64) *Response {
	return &Response{
		Code:    200,
		Message: "success",
		Data: map[string]interface{}{
			"list": data,
			"page": PageResponse{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
		},
	}
}

// Error 错误响应
func Error(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
	}
}
