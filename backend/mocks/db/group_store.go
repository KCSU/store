// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	model "github.com/kcsu/store/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// GroupStore is an autogenerated mock type for the GroupStore type
type GroupStore struct {
	mock.Mock
}

// AddUser provides a mock function with given fields: group, email
func (_m *GroupStore) AddUser(group *model.Group, email string) error {
	ret := _m.Called(group, email)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Group, string) error); ok {
		r0 = rf(group, email)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: group
func (_m *GroupStore) Create(group *model.Group) error {
	ret := _m.Called(group)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Group) error); ok {
		r0 = rf(group)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: group
func (_m *GroupStore) Delete(group *model.Group) error {
	ret := _m.Called(group)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Group) error); ok {
		r0 = rf(group)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: id
func (_m *GroupStore) Find(id uuid.UUID) (model.Group, error) {
	ret := _m.Called(id)

	var r0 model.Group
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.Group); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.Group)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields:
func (_m *GroupStore) Get() ([]model.Group, error) {
	ret := _m.Called()

	var r0 []model.Group
	if rf, ok := ret.Get(0).(func() []model.Group); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveUser provides a mock function with given fields: group, email
func (_m *GroupStore) RemoveUser(group *model.Group, email string) error {
	ret := _m.Called(group, email)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Group, string) error); ok {
		r0 = rf(group, email)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReplaceLookupUsers provides a mock function with given fields: group, users
func (_m *GroupStore) ReplaceLookupUsers(group *model.Group, users []model.GroupUser) error {
	ret := _m.Called(group, users)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Group, []model.GroupUser) error); ok {
		r0 = rf(group, users)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: group
func (_m *GroupStore) Update(group *model.Group) error {
	ret := _m.Called(group)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Group) error); ok {
		r0 = rf(group)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewGroupStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewGroupStore creates a new instance of GroupStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGroupStore(t mockConstructorTestingTNewGroupStore) *GroupStore {
	mock := &GroupStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
