package models

type CustomerProduct struct {
	Customer_id uint64  `gorm:"primaryKey;"`
	Product_id  uint64  `gorm:"primaryKey;"`
	Quantity    uint64  `json:"quantity"`
	TotalPrice  float32 `json:"total_price"`
}

type CustomersProducts struct {
	CustomerProductList []CustomerProduct `json:"customer_product"`
}
