// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	models "blinktag.com/bikesy-wrapper/models"
	mock "github.com/stretchr/testify/mock"
)

// ElevationService is an autogenerated mock type for the ElevationService type
type ElevationService struct {
	mock.Mock
}

// GetElevationsAndDistances provides a mock function with given fields: nodes, distances
func (_m *ElevationService) GetElevationsAndDistances(nodes []int, distances []float32) ([]models.ElevationDistance, error) {
	ret := _m.Called(nodes, distances)

	var r0 []models.ElevationDistance
	if rf, ok := ret.Get(0).(func([]int, []float32) []models.ElevationDistance); ok {
		r0 = rf(nodes, distances)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ElevationDistance)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]int, []float32) error); ok {
		r1 = rf(nodes, distances)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
