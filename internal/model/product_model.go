package model

import "github.com/google/uuid"

type GetProductRequest struct {
	Date string `json:"-" validate:"required,datetime=2006-01-02"`
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
