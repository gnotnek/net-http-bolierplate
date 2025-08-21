package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func New(redisURL string) *redis.Client {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal().Msgf("failed to parse redis url: %s", err.Error())
	}

	return redis.NewClient(opts)
}
