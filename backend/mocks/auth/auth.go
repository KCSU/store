// Code generated by mockery 2.10.6. DO NOT EDIT.

package mocks

import (
	echo "github.com/labstack/echo/v4"
	goth "github.com/markbates/goth"

	mock "github.com/stretchr/testify/mock"
)

// Auth is an autogenerated mock type for the Auth type
type Auth struct {
	mock.Mock
}

// CompleteUserAuth provides a mock function with given fields: c
func (_m *Auth) CompleteUserAuth(c echo.Context) (goth.User, error) {
	ret := _m.Called(c)

	var r0 goth.User
	if rf, ok := ret.Get(0).(func(echo.Context) goth.User); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Get(0).(goth.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(echo.Context) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAuthUrl provides a mock function with given fields: c
func (_m *Auth) GetAuthUrl(c echo.Context) (string, error) {
	ret := _m.Called(c)

	var r0 string
	if rf, ok := ret.Get(0).(func(echo.Context) string); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(echo.Context) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserId provides a mock function with given fields: c
func (_m *Auth) GetUserId(c echo.Context) int {
	ret := _m.Called(c)

	var r0 int
	if rf, ok := ret.Get(0).(func(echo.Context) int); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}
