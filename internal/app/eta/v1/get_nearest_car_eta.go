package v1

import (
	"context"

	"github.com/erizaver/eta-service/pkg/api"
)

func (e *EtaApi) GetNearestCarEtaV1(ctx context.Context, req *api.GetNearestCarEtaV1Request) (*api.GetNearestCarEtaV1Response, error) {
	nearestCar, err := e.EtaService.GetNearestCarEta(ctx, req.GetLatitude(), req.GetLongitude())
	if err != nil {
		return nil, err
	}

	return &api.GetNearestCarEtaV1Response{
		Data: &api.GetNearestCarEtaV1Response_Data{
			CarId: nearestCar.ID,
			Eta:   nearestCar.ArrivalMinutes,
		},
	}, nil
}
