package models

type Order struct {
	Order_id    uint64  `gorm:"primaryKey;autoIncrement"`
	Order_date  string  `json:"order_date"`
	Total_price float32 `json:"total_price"`
	Customer_id uint64  `json:"customer_id"`
	Payment_id  uint64  `json:"payment_id"`
	Shipment_id uint64  `json:"shipment_id"`
}
