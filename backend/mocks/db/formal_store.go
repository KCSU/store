// Code generated by mockery 2.12.1. DO NOT EDIT.

package mocks

import (
	model "github.com/kcsu/store/model"
	mock "github.com/stretchr/testify/mock"

	testing "testing"

	uuid "github.com/google/uuid"
)

// FormalStore is an autogenerated mock type for the FormalStore type
type FormalStore struct {
	mock.Mock
}

// All provides a mock function with given fields:
func (_m *FormalStore) All() ([]model.Formal, error) {
	ret := _m.Called()

	var r0 []model.Formal
	if rf, ok := ret.Get(0).(func() []model.Formal); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Formal)
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

// Create provides a mock function with given fields: formal
func (_m *FormalStore) Create(formal *model.Formal) error {
	ret := _m.Called(formal)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Formal) error); ok {
		r0 = rf(formal)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: formal
func (_m *FormalStore) Delete(formal *model.Formal) error {
	ret := _m.Called(formal)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Formal) error); ok {
		r0 = rf(formal)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: id
func (_m *FormalStore) Find(id uuid.UUID) (model.Formal, error) {
	ret := _m.Called(id)

	var r0 model.Formal
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.Formal); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.Formal)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindGuestList provides a mock function with given fields: id
func (_m *FormalStore) FindGuestList(id uuid.UUID) ([]model.FormalGuest, error) {
	ret := _m.Called(id)

	var r0 []model.FormalGuest
	if rf, ok := ret.Get(0).(func(uuid.UUID) []model.FormalGuest); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FormalGuest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindWithGroups provides a mock function with given fields: id
func (_m *FormalStore) FindWithGroups(id uuid.UUID) (model.Formal, error) {
	ret := _m.Called(id)

	var r0 model.Formal
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.Formal); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.Formal)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindWithTickets provides a mock function with given fields: id
func (_m *FormalStore) FindWithTickets(id uuid.UUID) (model.Formal, error) {
	ret := _m.Called(id)

	var r0 model.Formal
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.Formal); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.Formal)
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
func (_m *FormalStore) Get() ([]model.Formal, error) {
	ret := _m.Called()

	var r0 []model.Formal
	if rf, ok := ret.Get(0).(func() []model.Formal); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Formal)
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

// GetActive provides a mock function with given fields:
func (_m *FormalStore) GetActive() ([]model.Formal, error) {
	ret := _m.Called()

	var r0 []model.Formal
	if rf, ok := ret.Get(0).(func() []model.Formal); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Formal)
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

// GetGroups provides a mock function with given fields: ids
func (_m *FormalStore) GetGroups(ids []uuid.UUID) ([]model.Group, error) {
	ret := _m.Called(ids)

	var r0 []model.Group
	if rf, ok := ret.Get(0).(func([]uuid.UUID) []model.Group); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]uuid.UUID) error); ok {
		r1 = rf(ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWithUserData provides a mock function with given fields: userId
func (_m *FormalStore) GetWithUserData(userId uuid.UUID) ([]model.Formal, error) {
	ret := _m.Called(userId)

	var r0 []model.Formal
	if rf, ok := ret.Get(0).(func(uuid.UUID) []model.Formal); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Formal)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TicketsRemaining provides a mock function with given fields: formal, isGuest
func (_m *FormalStore) TicketsRemaining(formal *model.Formal, isGuest bool) uint {
	ret := _m.Called(formal, isGuest)

	var r0 uint
	if rf, ok := ret.Get(0).(func(*model.Formal, bool) uint); ok {
		r0 = rf(formal, isGuest)
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// Update provides a mock function with given fields: formal
func (_m *FormalStore) Update(formal *model.Formal) error {
	ret := _m.Called(formal)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Formal) error); ok {
		r0 = rf(formal)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateGroups provides a mock function with given fields: formal, groups
func (_m *FormalStore) UpdateGroups(formal model.Formal, groups []model.Group) error {
	ret := _m.Called(formal, groups)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.Formal, []model.Group) error); ok {
		r0 = rf(formal, groups)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewFormalStore creates a new instance of FormalStore. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewFormalStore(t testing.TB) *FormalStore {
	mock := &FormalStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
