package models

type Payment struct {
	Payment_id     uint64  `gorm:"primaryKey;autoIncrement"`
	Payment_date   string  `json:"payment_date"`
	Payment_method string  `json:"payment_method"`
	Amount         float32 `json:"amount"`
	Customer_id    uint64  `json:"customer_id"`
}
