package config

import "time"

type Config struct {
	Timeout                                time.Duration `envconfig:"HTTP_TIMEOUT_SECONDS"`
	RetryAmount                            int           `envconfig:"HTTP_RETRY_ATTEMPTS"`
	WheelyHost                             string        `envconfig:"WHEELY_FAKE_SERVICE_HOST"`
	WheelyFakeEtaBasePath                  string        `envconfig:"WHEELY_FAKE_ETA_BASE_PATH"`
	WheelyFakeCarServiceCarLimitPerRequest int64         `envconfig:"WHEELY_FAKE_CAR_SERVICE_CAR_LIMIT_PER_REQUEST"`
	HttpPort                               string        `envconfig:"HTTP_PORT"`
	RedisHost                              string        `envconfig:"REDIS_HOST"`
	RedisPort                              string        `envconfig:"REDIS_PORT"`
	RedisCacheInvalidationTime             time.Duration `envconfig:"REDIS_CACHE_INVALIDATION_TIME_SECONDS"`
	RedisCacheSearchRadiusMeters           float64       `envconfig:"REDIS_CACHE_SEARCH_RADIUS_METERS"`
}
