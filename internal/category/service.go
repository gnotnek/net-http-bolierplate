package category

import (
	"context"
	"net-http-boilerplate/internal/entity"
	apperror "net-http-boilerplate/internal/pkg/app-error"

	"gorm.io/gorm"
)

type Service struct {
	repo Repo
}

type Repo interface {
	Create(ctx context.Context, category *entity.Category) error
	FindAll(ctx context.Context, filter *entity.Filter) ([]entity.Category, *entity.Stats, error)
	FindByID(ctx context.Context, id int) (*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id int) error
}

func NewCategoryService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, req *CreateCategoryRequest) error {
	category := &entity.Category{
		Name: req.Name,
	}

	return s.repo.Create(ctx, category)
}

func (s *Service) FindAll(ctx context.Context, filter *entity.Filter) ([]CategoryResponse, *entity.Stats, error) {
	categories, stats, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, stats, apperror.ErrResourceNotFound
		}
		return nil, stats, err
	}

	var response []CategoryResponse
	for _, category := range categories {
		response = append(response, CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	return response, stats, nil
}

func (s *Service) FindByID(ctx context.Context, id int) (*entity.Category, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.ErrResourceNotFound
		}
	}

	return category, nil
}

func (s *Service) Update(ctx context.Context, id int, req UpdateCategoryRequest) (*CategoryResponse, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.ErrResourceNotFound
		}

		return nil, err
	}

	existing.Name = req.Name

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return &CategoryResponse{
		ID:   existing.ID,
		Name: existing.Name,
	}, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	cat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperror.ErrResourceNotFound
		}
	}

	return s.repo.Delete(ctx, cat.ID)
}
