package models

type Customer struct {
	Customer_id  uint64 `gorm:"primaryKey;autoIncrement"`
	First_name   string `json:"first_name"`
	Last_name    string `json:"last_name"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	Phone_Number string `json:"phone_number"`
}

type Customers struct {
	CustomerList []Customer `json:"customers"`
}

type CreateCustomerRequest struct {
	First_name   string `json:"first_name"`
	Last_name    string `json:"last_name"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	Phone_Number string `json:"phone_number"`
}
