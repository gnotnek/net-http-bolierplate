package post

import (
	"context"
	"net-http-boilerplate/internal/entity"
	apperror "net-http-boilerplate/internal/pkg/app-error"
	"strings"

	"gorm.io/gorm"
)

type Service struct {
	repo Repo
}

type Repo interface {
	Create(ctx context.Context, post *entity.Post) error
	FindAll(ctx context.Context, filter *entity.Filter) ([]entity.Post, *entity.Stats, error)
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

func (s *Service) Create(ctx context.Context, req *CreatePostRequest) (*entity.Post, error) {
	post := &entity.Post{
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
	}

	slug := strings.ReplaceAll(strings.ToLower(post.Title), " ", "-")
	post.Slug = slug
	if err := s.repo.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Service) FindAll(ctx context.Context, filter *entity.Filter) ([]PostResponse, *entity.Stats, error) {
	posts, stats, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, stats, apperror.ErrResourceNotFound
		}
		return nil, stats, err
	}

	var res []PostResponse
	for _, post := range posts {
		res = append(res, PostResponse{
			ID:      post.ID,
			Title:   post.Title,
			Content: post.Content,
			Slug:    post.Slug,
		})
	}

	return res, stats, nil
}

func (s *Service) FindByID(ctx context.Context, id int) (*PostResponse, error) {
	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.ErrResourceNotFound
		}
		return nil, err
	}
	return &PostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Slug:    post.Slug,
	}, nil
}

func (s *Service) Update(ctx context.Context, post *entity.Post) error {
	post.Slug = strings.ReplaceAll(strings.ToLower(post.Title), " ", "-")
	err := s.repo.Update(ctx, post)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperror.ErrResourceNotFound
		}
	}
	return err
}

func (s *Service) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperror.ErrResourceNotFound
		}
		return err
	}
	return nil
}
