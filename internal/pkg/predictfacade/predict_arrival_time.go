package predictfacade

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/erizaver/eta-service/internal/pkg/externalclients/predictservice/client/operations"
	"github.com/erizaver/eta-service/internal/pkg/externalclients/predictservice/models"
	"github.com/erizaver/eta-service/internal/pkg/model"
)

func (f *PredictFacade) PredictArrivalTime(ctx context.Context, cars []model.Car, lat, lng float64) ([]int64, error) {
	source := make([]models.Position, len(cars))
	for k, v := range cars {
		source[k] = models.Position{
			Lat: v.Lat,
			Lng: v.Lng,
		}
	}

	params := &operations.PredictParams{
		PositionList: operations.PredictBody{
			Source: source,
			Target: models.Position{
				Lat: lat,
				Lng: lng,
			},
		},
	}

	params.WithContext(ctx)

	eta, err := f.PredictClient.Operations.Predict(params)
	if err != nil {
		log.Errorf("unable to predict eta with wheely predict service, error:%s", err)
		return nil, errors.Wrap(err, "unable to predict time with external client")
	}

	return eta.Payload, nil
}
