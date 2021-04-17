package predictfacade

import (
	"github.com/erizaver/eta-service/internal/pkg/externalclients/predictservice/client"
)

type PredictFacade struct {
	PredictClient *client.PredictService
}

func NewPredictFacade(pc *client.PredictService) *PredictFacade {
	return &PredictFacade{
		PredictClient: pc,
	}
}
