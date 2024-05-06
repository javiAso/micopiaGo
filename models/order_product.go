package models

type OrderProduct struct {
	Order_id   uint64  `json:"order_id"`
	Product_id uint64  `json:"product_id"`
	Quantity   uint64  `json:"quantity"`
	Price      float32 `json:"total_price"`
}
