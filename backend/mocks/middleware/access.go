// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	echo "github.com/labstack/echo/v4"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kcsu/store/model"
)

// Access is an autogenerated mock type for the Access type
type Access struct {
	mock.Mock
}

// Get provides a mock function with given fields: page, size
func (_m *Access) Get(page int, size int) ([]model.AccessLog, error) {
	ret := _m.Called(page, size)

	var r0 []model.AccessLog
	if rf, ok := ret.Get(0).(func(int, int) []model.AccessLog); ok {
		r0 = rf(page, size)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.AccessLog)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(page, size)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Log provides a mock function with given fields: c, verb, metadata
func (_m *Access) Log(c echo.Context, verb string, metadata map[string]string) error {
	ret := _m.Called(c, verb, metadata)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, string, map[string]string) error); ok {
		r0 = rf(c, verb, metadata)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewAccess interface {
	mock.TestingT
	Cleanup(func())
}

// NewAccess creates a new instance of Access. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAccess(t mockConstructorTestingTNewAccess) *Access {
	mock := &Access{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
