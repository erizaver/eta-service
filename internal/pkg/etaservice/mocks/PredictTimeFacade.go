// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/erizaver/eta-service/internal/pkg/model"
)

// PredictTimeFacade is an autogenerated mock type for the PredictTimeFacade type
type PredictTimeFacade struct {
	mock.Mock
}

// PredictArrivalTime provides a mock function with given fields: ctx, cars, lat, lng
func (_m *PredictTimeFacade) PredictArrivalTime(ctx context.Context, cars []model.Car, lat float64, lng float64) ([]int64, error) {
	ret := _m.Called(ctx, cars, lat, lng)

	var r0 []int64
	if rf, ok := ret.Get(0).(func(context.Context, []model.Car, float64, float64) []int64); ok {
		r0 = rf(ctx, cars, lat, lng)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int64)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []model.Car, float64, float64) error); ok {
		r1 = rf(ctx, cars, lat, lng)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
