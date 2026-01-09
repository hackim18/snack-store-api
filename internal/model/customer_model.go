package model

type GetCustomerRequest struct {
	Page     int `json:"-" validate:"gte=1"`
	PageSize int `json:"-" validate:"gte=1"`
}

type CustomerResponse struct {
	Name   string `json:"name,omitempty"`
	Points int    `json:"points,omitempty"`
}
