package v1

import (
	"context"

	"github.com/erizaver/eta-service/internal/pkg/model"
	"github.com/erizaver/eta-service/pkg/api"
)

type EtaApi struct {
	EtaService EtaService
	api.UnimplementedEtaServiceServer
}

func NewEtaApi(es EtaService) *EtaApi {
	return &EtaApi{
		EtaService: es,
	}
}

type EtaService interface {
	GetNearestCarEta(ctx context.Context, lat, lng float64) (model.NearestCar, error)
}
