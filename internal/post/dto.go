package post

type CreatePostRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	ImageURL   string `json:"image_url"`
	CategoryID int    `json:"category_id"`
}

type UpdatePostRequest struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	ImageURL   string `json:"image_url"`
	CategoryID int    `json:"category_id"`
}
