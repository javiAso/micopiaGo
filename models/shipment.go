package models

type Shipment struct {
	Shipment_id   uint64 `gorm:"primaryKey;autoIncrement"`
	Shipment_date string `json:"shipment_date"`
	Address       string `json:"address"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	Zip_code      string `json:"zip_code"`
	Customer_id   uint64 `json:"customer_id"`
}
