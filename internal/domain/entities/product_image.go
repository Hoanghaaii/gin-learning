package entities

import "time"

type ProductImage struct {
	ID        uint
	ProductID uint
	URL       string
	IsCover   bool
	CreatedAt time.Time
}
