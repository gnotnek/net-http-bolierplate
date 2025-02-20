package category

import (
	"context"
	"net-http-boilerplate/internal/entity"

	"gorm.io/gorm"
)

type Service struct {
	repo Repo
}

type Repo interface {
	Create(ctx context.Context, category *entity.Category) error
	FindAll(ctx context.Context) ([]entity.Category, error)
	FindByID(ctx context.Context, id int) (*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id int) error
}

func NewCategoryService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, category *entity.Category) error {
	cat, err := s.repo.FindByID(ctx, category.ID)
	if err == nil && cat != nil {
		return ErrCategoryAlreadyExists
	}

	return s.repo.Create(ctx, category)
}

func (s *Service) FindAll(ctx context.Context) ([]entity.Category, error) {
	var categories []entity.Category
	categories, err := s.repo.FindAll(ctx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCategoryNotFound
		}
	}

	return categories, nil
}

func (s *Service) FindByID(ctx context.Context, id int) (*entity.Category, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCategoryNotFound
		}
	}

	return category, nil
}

func (s *Service) Update(ctx context.Context, category *entity.Category) error {
	cat, err := s.repo.FindByID(ctx, category.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrCategoryNotFound
		}
	}

	category.ID = cat.ID
	return s.repo.Update(ctx, category)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	cat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrCategoryNotFound
		}
	}

	return s.repo.Delete(ctx, cat.ID)
}
