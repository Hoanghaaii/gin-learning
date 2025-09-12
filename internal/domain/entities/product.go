package entities

import "time"

type ProductStatus string

const (
	ProductStatusActive     ProductStatus = "active"
	ProductStatusInActive   ProductStatus = "inactive"
	ProductStatusOutOfStock ProductStatus = "out_of_stock"
)

type Product struct {
	ID          uint
	SellerID    uint
	Name        string
	Description string
	BasePrice   float64
	Stock       uint
	Images      []ProductImage
	Status      ProductStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
