package models

type Product struct {
	Product_id  uint64  `gorm:"primaryKey;autoIncrement"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Stock       uint64  `json:"stock"`
	Category_id uint64  `json:"category_id"`
}

type Products struct {
	ProductList []Product `json:"products"`
}

type CreateProductRequest struct {
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Stock       uint64  `json:"stock"`
	Category_id uint64  `json:"category_id"`
}
