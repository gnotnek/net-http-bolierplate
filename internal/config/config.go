package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Database Database
	JWT      JWT
}

func Load() *Config {
	var c Config
	if err := env.Parse(&c); err != nil {
		log.Fatal(err)
	}

	return &c
}

type Database struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
}

type JWT struct {
	Secret string `env:"JWT_SECRET"`
	Issuer string `env:"JWT_ISSUER"`
}

func (d Database) DataSourceName() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.Name)
}
