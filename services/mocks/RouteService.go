// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	models "blinktag.com/bikesy-wrapper/models"
	mock "github.com/stretchr/testify/mock"
)

// RouteService is an autogenerated mock type for the RouteService type
type RouteService struct {
	mock.Mock
}

// GetBikeRoute provides a mock function with given fields: lat1, lng1, lat2, lng2
func (_m *RouteService) GetBikeRoute(lat1 string, lng1 string, lat2 string, lng2 string) (models.RouteResponse, error) {
	ret := _m.Called(lat1, lng1, lat2, lng2)

	var r0 models.RouteResponse
	if rf, ok := ret.Get(0).(func(string, string, string, string) models.RouteResponse); ok {
		r0 = rf(lat1, lng1, lat2, lng2)
	} else {
		r0 = ret.Get(0).(models.RouteResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string) error); ok {
		r1 = rf(lat1, lng1, lat2, lng2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetProfile provides a mock function with given fields: profile
func (_m *RouteService) SetProfile(profile string) {
	_m.Called(profile)
}
