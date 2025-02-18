package entity

type Category struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}
