package post

import (
	"context"
	"net-http-boilerplate/internal/entity"

	"gorm.io/gorm"
)

type Service struct {
	repo Repo
}

type Repo interface {
	Create(ctx context.Context, post *entity.Post) error
	FindAll(ctx context.Context) ([]entity.Post, error)
	FindByCategory(ctx context.Context, category string) ([]entity.Post, error)
	FindByID(ctx context.Context, id int) (*entity.Post, error)
	Update(ctx context.Context, post *entity.Post) error
	Delete(ctx context.Context, id int) error
}

func NewPostService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, post *entity.Post) error {
	return s.repo.Create(ctx, post)
}

func (s *Service) FindAll(ctx context.Context) ([]entity.Post, error) {
	var posts []entity.Post
	posts, err := s.repo.FindAll(ctx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	return posts, nil
}

func (s *Service) FindByCategory(ctx context.Context, category string) ([]entity.Post, error) {
	var posts []entity.Post
	posts, err := s.repo.FindByCategory(ctx, category)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	return posts, nil
}

func (s *Service) FindByID(ctx context.Context, id int) (*entity.Post, error) {
	var post *entity.Post
	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	return post, nil
}

func (s *Service) Update(ctx context.Context, post *entity.Post) error {
	return s.repo.Update(ctx, post)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrPostNotFound
		}
		return err
	}
	return nil
}
