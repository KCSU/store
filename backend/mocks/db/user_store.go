// Code generated by mockery 2.10.0. DO NOT EDIT.

package mocks

import (
	goth "github.com/markbates/goth"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kcsu/store/model"
)

// UserStore is an autogenerated mock type for the UserStore type
type UserStore struct {
	mock.Mock
}

// Exists provides a mock function with given fields: email
func (_m *UserStore) Exists(email string) (bool, error) {
	ret := _m.Called(email)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: id
func (_m *UserStore) Find(id int) (model.User, error) {
	ret := _m.Called(id)

	var r0 model.User
	if rf, ok := ret.Get(0).(func(int) model.User); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOrCreate provides a mock function with given fields: gu
func (_m *UserStore) FindOrCreate(gu *goth.User) (model.User, error) {
	ret := _m.Called(gu)

	var r0 model.User
	if rf, ok := ret.Get(0).(func(*goth.User) model.User); ok {
		r0 = rf(gu)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*goth.User) error); ok {
		r1 = rf(gu)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Groups provides a mock function with given fields: user
func (_m *UserStore) Groups(user *model.User) ([]model.Group, error) {
	ret := _m.Called(user)

	var r0 []model.Group
	if rf, ok := ret.Get(0).(func(*model.User) []model.Group); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Permissions provides a mock function with given fields: user
func (_m *UserStore) Permissions(user *model.User) ([]model.Permission, error) {
	ret := _m.Called(user)

	var r0 []model.Permission
	if rf, ok := ret.Get(0).(func(*model.User) []model.Permission); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Permission)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
