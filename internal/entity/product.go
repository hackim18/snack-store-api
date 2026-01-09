package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	SizeSmall  = "Small"
	SizeMedium = "Medium"
	SizeLarge  = "Large"
)

func PointsCost(size string) int {
	switch strings.TrimSpace(size) {
	case SizeSmall:
		return 200
	case SizeMedium:
		return 300
	case SizeLarge:
		return 500
	default:
		return 0
	}
}

type Product struct {
	ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name             string    `gorm:"not null;check:length(btrim(name)) > 0"`
	Type             string    `gorm:"column:type;not null;check:length(btrim(type)) > 0;index:products_type_idx"`
	Flavor           string    `gorm:"not null;check:flavor IN ('Jagung Bakar','Rumput Laut','Original','Jagung Manis','Keju Asin','Keju Manis','Pedas');index:products_flavor_idx"`
	Size             string    `gorm:"type:varchar(10);not null;check:size IN ('Small','Medium','Large');index:products_size_idx"`
	Price            int       `gorm:"not null;check:price >= 0"`
	StockQty         int       `gorm:"column:stock_qty;not null;check:stock_qty >= 0"`
	ManufacturedDate time.Time `gorm:"type:date;not null;index:products_manufactured_date_idx"`
	CreatedAt        time.Time `gorm:"not null;default:now()"`
	UpdatedAt        time.Time `gorm:"not null;default:now()"`
}

func (p *Product) TableName() string {
	return "products"
}

func (p *Product) BeforeCreate(_ *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	return
}
