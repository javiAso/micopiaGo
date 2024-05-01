package models

type CustomerProduct struct {
	Customer_product_id uint64  `gorm:"primaryKey;autoIncrement"`
	Customer_id         uint64  `json:"customer_id"`
	Product_id          uint64  `json:"product_id"`
	Quantity            uint64  `json:"quantity"`
	TotalPrice          float32 `json:"total_price"`
}

type CustomersProducts struct {
	CustomerProductList []CustomerProduct `json:"customer_product"`
}
