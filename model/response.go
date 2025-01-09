package model

import (
	"math"
)

// BaseResponse 基础响应结构
type BaseResponse struct {
	Success bool   `json:"success"` // 响应成功
	Message string `json:"message"` // 响应信息
}

// PageInfo 分页信息
type PageInfo struct {
	PageIndex  int   `json:"pageIndex"`  // 当前页码
	PageSize   int   `json:"pageSize"`   // 每页大小
	Total      int64 `json:"total"`      // 总记录数
	TotalPages int   `json:"totalPages"` // 总页数
}

// PageResponse 分页响应结构
type PageResponse struct {
	Success bool        `json:"success"` // 响应成功
	Message string      `json:"message"` // 响应信息
	Data    interface{} `json:"data"`    // 列表数据
	Page    PageInfo    `json:"page"`    // 分页信息
}

// ListResponse 列表响应结构（不分页）
type ListResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
	Total   int64         `json:"total"`
}

// EmptyResponse 空响应结构
type EmptyResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// DetailResponse 详情响应结构
type DetailResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// BatchResponse 批量操作响应结构
type BatchResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	SuccessList []uint `json:"successList,omitempty"` // 成功的ID列表
	FailedList  []uint `json:"failedList,omitempty"`  // 失败的ID列表
}

// StatusResponse 状态响应结构
type StatusResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// NewResponse 创建基础响应
func NewResponse(success bool, message string) BaseResponse {
	return BaseResponse{
		Success: success,
		Message: message,
	}
}

// NewPageResponse 创建分页响应
func NewPageResponse(success bool, message string, data interface{}, page, pageSize int, total int64) PageResponse {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return PageResponse{
		Success: success,
		Message: message,
		Data:    data,
		Page: PageInfo{
			PageIndex:  page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}

// NewListResponse 创建列表响应
func NewListResponse(success bool, message string, data []interface{}, total int64) ListResponse {
	return ListResponse{
		Success: success,
		Message: message,
		Data:    data,
		Total:   total,
	}
}

// NewBatchResponse 创建批量操作响应
func NewBatchResponse(success bool, message string, successList []uint, failedList []uint) BatchResponse {
	return BatchResponse{
		Success:     success,
		Message:     message,
		SuccessList: successList,
		FailedList:  failedList,
	}
}

// NewDetailResponse 创建详情响应
func NewDetailResponse(success bool, message string, data interface{}) DetailResponse {
	return DetailResponse{
		Success: success,
		Message: message,
		Data:    data,
	}
}
