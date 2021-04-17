package redisfacade

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/erizaver/eta-service/internal/pkg/model"
)

func (f *RedisFacade) AddEtaToCache(car model.NearestCar, lat, lng float64) error {
	dataToCache, err := json.Marshal(model.CachedCar{
		ID:             car.ID,
		ArrivalMinutes: car.ArrivalMinutes,
		TimeCreated:    time.Now(),
	})
	if err != nil {
		log.Errorf("unable to marshal data, carID: %d , eta: %d, err:%s", car.ID, car.ArrivalMinutes, err)
		return errors.Wrap(err, "unable to marshal data")
	}

	res := f.RedisClient.GeoAdd(etaRedisKey, &redis.GeoLocation{
		Name:      string(dataToCache),
		Longitude: lng,
		Latitude:  lat,
	})

	if res.Err() != nil {
		log.Errorf("unable to add data to cache, error: %s", res.Err())
		return errors.Wrap(res.Err(), "unable to add data to cache")
	}

	return nil
}
