package category

import (
	"context"
	"net-http-boilerplate/internal/entity"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *Repository) FindAll(ctx context.Context) ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.WithContext(ctx).Find(&categories).Error
	return categories, err
}

func (r *Repository) FindByID(ctx context.Context, id int) (*entity.Category, error) {
	var category entity.Category
	err := r.db.WithContext(ctx).First(&category, id).Error
	return &category, err
}

func (r *Repository) Update(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&entity.Category{}, id).Error
}
