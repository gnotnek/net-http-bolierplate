package postgres

import (
	"net-http-boilerplate/internal/config"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxIdleConns    = 10
	maxOpenConns    = 100
	connMaxLifetime = time.Hour
)

// NewGORM initializes a new GORM database connection using the provided configuration.
func NewGORM(c *config.Database) *gorm.DB {
	db, err := gorm.Open(postgres.Open(c.DataSourceName()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database connection")
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	return db
}
