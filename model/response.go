package model

type Response struct {
	Success bool      `json:"success,omitempty"`
	Message string    `json:"message,omitempty"`
	Data    []Content `json:"data,omitempty"`
	Page    *Page     `json:"pagination,omitempty"`
}

type TestResponse struct {
	Success bool        `json:"success,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Page    *Page       `json:"page,omitempty"`
}

type Page struct {
	PageSize     int    `json:"pageSize,omitempty"`
	PageIndex    int    `json:"pageIndex,omitempty"`
	ContentCount int64  `json:"total,omitempty"`
	PageCount    int64  `json:"pageCount,omitempty"`
	Next         string `json:"next,omitempty"`
	Pre          string `json:"pre,omitempty"`
	Order        string `json:"orderBy,omitempty"`
}
