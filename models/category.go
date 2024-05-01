package models

type Category struct {
	Category_id uint64 `gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name"`
}

type Categories struct {
	CategoryList []Category `json:"categories"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}
