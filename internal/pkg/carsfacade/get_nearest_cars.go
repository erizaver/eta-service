package carsfacade

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/erizaver/eta-service/internal/pkg/externalclients/carservice/client/operations"
	"github.com/erizaver/eta-service/internal/pkg/model"
)

func (f *CarsFacade) GetNearestCars(ctx context.Context, lat, lng float64) ([]model.Car, error) {
	params := &operations.GetCarsParams{
		Lat:   lat,
		Limit: f.CarsLimitPerRequest,
		Lng:   lng,
	}
	params.WithContext(ctx)

	cars, err := f.CarsClient.Operations.GetCars(params)
	if err != nil {
		log.Errorf("unable to get cars from wheely car service, error:%s", err)
		return nil, errors.Wrap(err, "unable to get cars from external service")
	}

	res := make([]model.Car, len(cars.GetPayload()))
	for k, v := range cars.GetPayload() {
		res[k] = model.Car{
			ID:  v.ID,
			Lat: v.Lat,
			Lng: v.Lng,
		}
	}

	return res, nil
}
