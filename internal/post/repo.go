package post

import (
	"context"
	"net-http-boilerplate/internal/entity"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, post *entity.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *Repository) FindAll(ctx context.Context) ([]entity.Post, error) {
	var posts []entity.Post
	err := r.db.WithContext(ctx).Find(&posts).Error
	return posts, err
}

func (r *Repository) FindByCategory(ctx context.Context, category string) ([]entity.Post, error) {
	var posts []entity.Post
	query := `
		SELECT * FROM posts
		JOIN categories ON posts.category_id = categories.id
		WHERE categories.name = ?
	`

	err := r.db.
		WithContext(ctx).
		Raw(query, category).
		Scan(&posts).
		Error
	return posts, err
}

func (r *Repository) FindByID(ctx context.Context, id int) (*entity.Post, error) {
	var post entity.Post
	err := r.db.WithContext(ctx).First(&post, id).Error
	return &post, err
}

func (r *Repository) Update(ctx context.Context, post *entity.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&entity.Post{}, id).Error
}
