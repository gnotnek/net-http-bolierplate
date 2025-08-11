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

func (r *Repository) FindAll(ctx context.Context, filter *entity.Filter) ([]entity.Post, *entity.Stats, error) {
	var res []entity.Post
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Post{})

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if filter.StartDate != nil && filter.EndDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate)
	} else if filter.StartDate != nil {
		query = query.Where("created_at >= ?", filter.StartDate)
	} else if filter.EndDate != nil {
		query = query.Where("created_at <= ?", filter.EndDate)
	}

	// Count total items
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	if filter.SortBy != nil && *filter.SortBy != "" {
		sortOrder := "ASC"
		if filter.SortOrder != nil && (*filter.SortOrder == "ASC" || *filter.SortOrder == "DESC") {
			sortOrder = *filter.SortOrder
		}
		query = query.Order(*filter.SortBy + " " + sortOrder)
	} else {
		query = query.Order("created_at ASC")
	}

	// Pagination logic
	if filter.Page != nil && filter.PerPage != nil {
		page := *filter.Page
		perPage := *filter.PerPage
		offset := (page - 1) * perPage

		if err := query.Limit(perPage).Offset(offset).Find(&res).Error; err != nil {
			return nil, nil, err
		}

		stats := &entity.Stats{
			Page:  page,
			Total: int(total),
			Limit: perPage,
		}
		return res, stats, nil
	}

	// If no pagination, return all data
	if err := query.Find(&res).Error; err != nil {
		return nil, nil, err
	}

	return res, nil, nil

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
