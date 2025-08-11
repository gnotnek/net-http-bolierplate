package entity

type Post struct {
	ID         int      `json:"id" gorm:"primaryKey"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Slug       string   `json:"slug"`
	AuthorID   int      `json:"author_id"`
	CategoryID int      `json:"category_id"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID;references:ID"`
}
