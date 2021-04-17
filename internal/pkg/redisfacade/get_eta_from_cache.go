package redisfacade

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/erizaver/eta-service/internal/pkg/errorhelper"
	"github.com/erizaver/eta-service/internal/pkg/model"
)

const etaRedisKey = "eta"

func (f *RedisFacade) GetEtaFromCache(lat, lng float64) (model.NearestCar, error) {
	res := model.NearestCar{}
	data := f.RedisClient.GeoRadius(etaRedisKey, lng, lat, &redis.GeoRadiusQuery{
		Radius: f.CacheSearchRadius,
		Unit:   "m",
	})

	if data.Err() != nil {
		log.Errorf("unable to get data from cache, error: %s", data.Err())
		return res, errors.Wrap(data.Err(), "unable to get data from cache")
	}

	if len(data.Val()) == 0 {
		return res, errorhelper.NewCacheMissError()
	}

	success := false
	for _, v := range data.Val() {
		car := &model.CachedCar{}
		err := json.Unmarshal([]byte(v.Name), car)
		if err != nil {
			log.Errorf("unable to unmarshal data from cache, error: %s data: %s", err, v.Name)
			return res, errors.Wrap(err, "unable to unmarshal data from cache")
		}

		if car.TimeCreated.Before(time.Now().Add(-f.CacheInvalidationTime)) {
			f.RedisClient.ZRem(etaRedisKey, v.Name)
			continue
		}

		if res.ID == 0 || (res.ID != 0 && res.ArrivalMinutes > car.ArrivalMinutes) {
			res.ID = car.ID
			res.ArrivalMinutes = car.ArrivalMinutes
			success = true
		}
	}

	if !success {
		return res, errorhelper.NewCacheMissError()
	}

	return res, nil
}
