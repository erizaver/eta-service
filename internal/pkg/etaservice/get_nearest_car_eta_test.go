package etaservice

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/erizaver/eta-service/internal/pkg/errorhelper"
	"github.com/erizaver/eta-service/internal/pkg/etaservice/mocks"
	"github.com/erizaver/eta-service/internal/pkg/model"
)

func TestGetNearestCarEta(t *testing.T) {
	t.Run("Successful call", func(t *testing.T) {
		carsFacadeMock := new(mocks.CarsFacade)
		carsFacadeMock.On("GetNearestCars",
			mock.Anything,
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64")).
			Once().Return([]model.Car{
			{
				ID:  1,
				Lat: 1.1,
				Lng: 1.1,
			},
			{
				ID:  2,
				Lat: 2.2,
				Lng: 2.2,
			},
			{
				ID:  3,
				Lat: 3.3,
				Lng: 3.3,
			},
		}, nil)

		predictFacadeMock := new(mocks.PredictTimeFacade)
		predictFacadeMock.On("PredictArrivalTime",
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64")).
			Once().Return([]int64{
			15,
			4,
			32,
		}, nil)

		redisFacade := new(mocks.RedisFacade)
		redisFacade.On("GetEtaFromCache",
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64")).Once().Return(
			model.NearestCar{}, errorhelper.NewCacheMissError())

		redisFacade.On("AddEtaToCache",
			mock.Anything,
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64")).Once().Return(nil)

		srv := NewEtaService(carsFacadeMock, predictFacadeMock, redisFacade)
		ctx := context.Background()

		res, err := srv.GetNearestCarEta(ctx, 50.50, 30.30)
		assert.NoError(t, err)
		assert.Equal(t, res.ID, int64(2))
		assert.Equal(t, res.ArrivalMinutes, int64(4))
	})

	t.Run("inconsistent data", func(t *testing.T) {
		carsFacadeMock := new(mocks.CarsFacade)
		carsFacadeMock.On("GetNearestCars",
			mock.Anything,
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64")).
			Once().Return([]model.Car{
			{
				ID:  1,
				Lat: 1.1,
				Lng: 1.1,
			},
			{
				ID:  2,
				Lat: 2.2,
				Lng: 2.2,
			},
		}, nil)

		predictFacadeMock := new(mocks.PredictTimeFacade)
		predictFacadeMock.On("PredictArrivalTime",
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64")).
			Once().Return([]int64{
			15,
			32,
			33,
		}, nil)

		redisFacade := new(mocks.RedisFacade)
		redisFacade.On("GetEtaFromCache",
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64")).Once().Return(
			model.NearestCar{}, errorhelper.NewCacheMissError())

		srv := NewEtaService(carsFacadeMock, predictFacadeMock, redisFacade)
		ctx := context.Background()

		_, err := srv.GetNearestCarEta(ctx, 50.50, 30.30)
		assert.Error(t, err)
	})

	t.Run("error from cars facade", func(t *testing.T) {
		carsFacadeMock := new(mocks.CarsFacade)
		carsFacadeMock.On("GetNearestCars",
			mock.Anything,
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64")).
			Once().Return(nil, errors.New("test error"))

		redisFacade := new(mocks.RedisFacade)
		redisFacade.On("GetEtaFromCache",
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64")).Once().Return(
			model.NearestCar{}, errorhelper.NewCacheMissError())

		srv := NewEtaService(carsFacadeMock, nil, redisFacade)
		ctx := context.Background()

		_, err := srv.GetNearestCarEta(ctx, 50.50, 30.30)
		assert.Error(t, err)
	})

	t.Run("error from predict facade", func(t *testing.T) {
		carsFacadeMock := new(mocks.CarsFacade)
		carsFacadeMock.On("GetNearestCars",
			mock.Anything,
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64")).
			Once().Return([]model.Car{
			{
				ID:  1,
				Lat: 1.1,
				Lng: 1.1,
			},
			{
				ID:  2,
				Lat: 2.2,
				Lng: 2.2,
			},
		}, nil)

		predictFacadeMock := new(mocks.PredictTimeFacade)
		predictFacadeMock.On("PredictArrivalTime",
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64")).
			Once().Return(nil, errors.New("test error"))

		redisFacade := new(mocks.RedisFacade)
		redisFacade.On("GetEtaFromCache",
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64")).Once().Return(
			model.NearestCar{}, errorhelper.NewCacheMissError())

		srv := NewEtaService(carsFacadeMock, predictFacadeMock, redisFacade)
		ctx := context.Background()

		_, err := srv.GetNearestCarEta(ctx, 50.50, 30.30)
		assert.Error(t, err)
	})

	t.Run("Cache hit", func(t *testing.T) {
		redisFacade := new(mocks.RedisFacade)
		redisFacade.On("GetEtaFromCache",
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64")).Once().Return(
			model.NearestCar{
				ID:             1,
				ArrivalMinutes: 2,
			}, nil)

		srv := NewEtaService(nil, nil, redisFacade)
		ctx := context.Background()

		res, err := srv.GetNearestCarEta(ctx, 50.50, 30.30)
		assert.NoError(t, err)
		assert.Equal(t, res.ID, int64(1))
		assert.Equal(t, res.ArrivalMinutes, int64(2))
	})

}
