package user

import (
	"context"
	"net-http-boilerplate/internal/entity"
	apperror "net-http-boilerplate/internal/pkg/app-error"
	"net-http-boilerplate/internal/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo Repo
	jwt  *jwt.JWT
}

type Repo interface {
	Create(ctx context.Context, user *entity.User) error
	Save(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

func NewUserService(repo Repo, jwt *jwt.JWT) *Service {
	return &Service{
		repo: repo,
		jwt:  jwt,
	}
}

func (s *Service) Register(ctx context.Context, req *RegisterRequest) error {
	user := &entity.User{
		Name:  req.Name,
		Email: req.Email,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return s.repo.Create(ctx, user)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*UserResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.ErrResourceNotFound
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperror.ErrInvalidPassword
	}

	token, _, err := s.jwt.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		AccessToken: token,
		Email:       user.Email,
		Name:        user.Name,
	}, nil
}

func (s *Service) Save(ctx context.Context, user *entity.User) error {
	return s.repo.Save(ctx, user)
}
