package model

type WebResponse[T any] struct {
	Message string        `json:"message,omitempty"`
	Data    T             `json:"data,omitempty"`
	Paging  *PageMetadata `json:"paging,omitempty"`
	Errors  string        `json:"errors,omitempty"`
}

type PageMetadata struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItem   int64 `json:"total_item"`
	TotalPage   int64 `json:"total_page"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}
