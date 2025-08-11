package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Posts     []Post    `json:"posts" gorm:"foreignKey:AuthorID;references:ID;constraint:OnDelete:SET NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
