package entity

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Posts    []Post `json:"posts" gorm:"foreignKey:AuthorID"`
}
