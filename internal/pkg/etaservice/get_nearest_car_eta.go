package etaservice

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/erizaver/eta-service/internal/pkg/errorhelper"
	"github.com/erizaver/eta-service/internal/pkg/model"
)

func (s *EtaService) GetNearestCarEta(ctx context.Context, lat, lng float64) (model.NearestCar, error) {
	carFromCache, err := s.RedisFacade.GetEtaFromCache(lat, lng)
	if err != nil {
		if !errorhelper.IfCacheMiss(err) {
			return model.NearestCar{}, errors.Wrap(err, "cache error")
		}
	} else {
		return carFromCache, nil
	}

	res := model.NearestCar{}

	nearestCars, err := s.CarsFacade.GetNearestCars(ctx, lat, lng)
	if err != nil {
		return res, errorhelper.WrapWithCode(err, errorhelper.Internal, "unable to get nearest cars")
	}

	if len(nearestCars) == 0 {
		return res, errorhelper.NewWithCode(errorhelper.NotFound, "no cars found nearby")
	}

	carsEta, err := s.PredictTimeFacade.PredictArrivalTime(ctx, nearestCars, lat, lng)
	if err != nil {
		return res, errorhelper.WrapWithCode(err, errorhelper.Internal, "unable to get eta")
	}

	if len(nearestCars) != len(carsEta) {
		log.Errorf("unable to get eta, inconsistent data, cars amount:%d, etas amount:%d", len(nearestCars), len(carsEta))
		return res, errorhelper.NewWithCode(errorhelper.BadRequest, "data is not valid, amount of cars is not the same")
	}

	bestTimeIndex := 0
	for k, v := range carsEta {
		if v < carsEta[bestTimeIndex] {
			bestTimeIndex = k
		}
	}

	// я не уверен, что нам нужен айдишник машины, этого нет в задании, но это звучит здраво,
	// ведь нам надо потом показать ее в приложении плюс вы же возвращаете его в фейк сервисе. :)
	// Лучше пусть будет, удалить не трудно.
	res.ID = nearestCars[bestTimeIndex].ID
	res.ArrivalMinutes = carsEta[bestTimeIndex]

	err = s.RedisFacade.AddEtaToCache(res, lat, lng)
	if err != nil {
		log.Errorf("unable to add data to cache, err:%s", err)
	}

	return res, nil
}
