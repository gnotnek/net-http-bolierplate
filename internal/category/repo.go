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

func (r *Repository) FindAll(ctx context.Context, filter *entity.Filter) ([]entity.Category, *entity.Stats, error) {
	var res []entity.Category
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.Category{})

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
