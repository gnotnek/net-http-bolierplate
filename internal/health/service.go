package health

import (
	"context"
	"net-http-boilerplate/internal/entity"

	"github.com/rs/zerolog/log"
)

type Service struct {
	repo *repo
}

func NewHealthService(repo *repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CheckDatabase(ctx context.Context) (*entity.HealthCheck, bool) {
	healthComponent := &entity.HealthCheck{
		Database: entity.HealthStatusOK,
	}

	if err := s.repo.CheckDatabase(ctx); err != nil {
		log.Ctx(ctx).Error().Msgf("Check database error: %s", err.Error())
		healthComponent.Database = entity.HealthStatusFailed
	}

	isHealthy := healthComponent.Database == entity.HealthStatusOK
	return healthComponent, isHealthy
}
