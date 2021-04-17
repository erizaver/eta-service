package etaservice

import (
	"context"

	"github.com/erizaver/eta-service/internal/pkg/model"
)

type EtaService struct {
	CarsFacade        CarsFacade
	PredictTimeFacade PredictTimeFacade
	RedisFacade       RedisFacade
}

func NewEtaService(cf CarsFacade, ptf PredictTimeFacade, rc RedisFacade) *EtaService {
	return &EtaService{
		CarsFacade:        cf,
		PredictTimeFacade: ptf,
		RedisFacade:       rc,
	}
}

type CarsFacade interface {
	GetNearestCars(ctx context.Context, lat, lng float64) ([]model.Car, error)
}

type PredictTimeFacade interface {
	PredictArrivalTime(ctx context.Context, cars []model.Car, lat, lng float64) ([]int64, error)
}

type RedisFacade interface {
	AddEtaToCache(car model.NearestCar, lat, lng float64) error
	GetEtaFromCache(lat, lng float64) (model.NearestCar, error)
}
