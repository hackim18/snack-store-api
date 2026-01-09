package model

import "github.com/google/uuid"

type GetProductRequest struct {
	Date string `json:"-" validate:"required,datetime=2006-01-02"`
}

type CreateProductRequest struct {
	Name             string `json:"name" validate:"required"`
	Type             string `json:"type" validate:"required"`
	Flavor           string `json:"flavor" validate:"required"`
	Size             string `json:"size" validate:"required,oneof=Small Medium Large"`
	Price            int    `json:"price" validate:"required,gte=0"`
	StockQty         int    `json:"stock_qty" validate:"required,gte=0"`
	ManufacturedDate string `json:"manufactured_date" validate:"required,datetime=2006-01-02"`
}

type ProductResponse struct {
	ID               *uuid.UUID `json:"id,omitempty"`
	Name             string     `json:"name,omitempty"`
	Type             string     `json:"type,omitempty"`
	Flavor           string     `json:"flavor,omitempty"`
	Size             string     `json:"size,omitempty"`
	Price            int        `json:"price,omitempty"`
	StockQty         int        `json:"stock_qty,omitempty"`
	ManufacturedDate string     `json:"manufactured_date,omitempty"`
}
