package entity

type Post struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	ImageURL string `json:"image_url"`
	Slug     string `json:"slug"`
	AuthorID int    `json:"author_id"`
}
