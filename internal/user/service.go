package user

import (
	"context"
	"net-http-boilerplate/internal/entity"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo Repo
}

type Repo interface {
	Create(ctx context.Context, user *entity.User) error
	Save(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

func NewUserService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(ctx context.Context, user *entity.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return s.repo.Create(ctx, user)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*entity.User, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

func (s *Service) Save(ctx context.Context, user *entity.User) error {
	return s.repo.Save(ctx, user)
}
