package postgres

import (
	"net-http-boilerplate/internal/entity"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&entity.HealthCheck{},
		&entity.User{},
		&entity.Category{},
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to auto migrate, err: %v", err.Error())
	}

	log.Info().Msg("migration completed")
}
