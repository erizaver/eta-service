package carsfacade

import (
	"github.com/erizaver/eta-service/internal/pkg/externalclients/carservice/client"
)

type CarsFacade struct {
	CarsClient          *client.CarsService
	CarsLimitPerRequest int64
}

func NewCarsFacade(cc *client.CarsService, cl int64) *CarsFacade {
	return &CarsFacade{
		CarsClient:          cc,
		CarsLimitPerRequest: cl,
	}
}
