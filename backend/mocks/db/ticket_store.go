// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	dto "github.com/kcsu/store/model/dto"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kcsu/store/model"

	uuid "github.com/google/uuid"
)

// TicketStore is an autogenerated mock type for the TicketStore type
type TicketStore struct {
	mock.Mock
}

// BatchCreate provides a mock function with given fields: tickets
func (_m *TicketStore) BatchCreate(tickets []model.Ticket) error {
	ret := _m.Called(tickets)

	var r0 error
	if rf, ok := ret.Get(0).(func([]model.Ticket) error); ok {
		r0 = rf(tickets)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CountGuestByFormal provides a mock function with given fields: formalId, userId
func (_m *TicketStore) CountGuestByFormal(formalId uuid.UUID, userId uuid.UUID) (int64, error) {
	ret := _m.Called(formalId, userId)

	var r0 int64
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) int64); ok {
		r0 = rf(formalId, userId)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(formalId, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ticket
func (_m *TicketStore) Create(ticket *model.Ticket) error {
	ret := _m.Called(ticket)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Ticket) error); ok {
		r0 = rf(ticket)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: id
func (_m *TicketStore) Delete(id uuid.UUID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteByFormal provides a mock function with given fields: formalID, userID
func (_m *TicketStore) DeleteByFormal(formalID uuid.UUID, userID uuid.UUID) error {
	ret := _m.Called(formalID, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(formalID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ExistsByFormal provides a mock function with given fields: formalID, userID
func (_m *TicketStore) ExistsByFormal(formalID uuid.UUID, userID uuid.UUID) (bool, error) {
	ret := _m.Called(formalID, userID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) bool); ok {
		r0 = rf(formalID, userID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(formalID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: id
func (_m *TicketStore) Find(id uuid.UUID) (model.Ticket, error) {
	ret := _m.Called(id)

	var r0 model.Ticket
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.Ticket); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.Ticket)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindWithFormal provides a mock function with given fields: id
func (_m *TicketStore) FindWithFormal(id uuid.UUID) (model.Ticket, error) {
	ret := _m.Called(id)

	var r0 model.Ticket
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.Ticket); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.Ticket)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: userId
func (_m *TicketStore) Get(userId uuid.UUID) ([]model.Ticket, error) {
	ret := _m.Called(userId)

	var r0 []model.Ticket
	if rf, ok := ret.Get(0).(func(uuid.UUID) []model.Ticket); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Ticket)
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

// Update provides a mock function with given fields: id, ticket
func (_m *TicketStore) Update(id uuid.UUID, ticket *dto.TicketRequestDto) error {
	ret := _m.Called(id, ticket)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, *dto.TicketRequestDto) error); ok {
		r0 = rf(id, ticket)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTicketStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewTicketStore creates a new instance of TicketStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTicketStore(t mockConstructorTestingTNewTicketStore) *TicketStore {
	mock := &TicketStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
