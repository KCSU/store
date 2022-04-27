// Code generated by mockery 2.12.1. DO NOT EDIT.

package mocks

import (
	model "github.com/kcsu/store/model"
	mock "github.com/stretchr/testify/mock"

	testing "testing"

	uuid "github.com/google/uuid"
)

// BillStore is an autogenerated mock type for the BillStore type
type BillStore struct {
	mock.Mock
}

// AddFormals provides a mock function with given fields: bill, formalIds
func (_m *BillStore) AddFormals(bill *model.Bill, formalIds []uuid.UUID) error {
	ret := _m.Called(bill, formalIds)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Bill, []uuid.UUID) error); ok {
		r0 = rf(bill, formalIds)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: id
func (_m *BillStore) Find(id uuid.UUID) (model.Bill, error) {
	ret := _m.Called(id)

	var r0 model.Bill
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.Bill); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.Bill)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindWithFormals provides a mock function with given fields: id
func (_m *BillStore) FindWithFormals(id uuid.UUID) (model.Bill, error) {
	ret := _m.Called(id)

	var r0 model.Bill
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.Bill); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.Bill)
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
func (_m *BillStore) Get() ([]model.Bill, error) {
	ret := _m.Called()

	var r0 []model.Bill
	if rf, ok := ret.Get(0).(func() []model.Bill); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Bill)
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

// GetCostBreakdown provides a mock function with given fields: bill
func (_m *BillStore) GetCostBreakdown(bill *model.Bill) ([]model.FormalCostBreakdown, error) {
	ret := _m.Called(bill)

	var r0 []model.FormalCostBreakdown
	if rf, ok := ret.Get(0).(func(*model.Bill) []model.FormalCostBreakdown); ok {
		r0 = rf(bill)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FormalCostBreakdown)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Bill) error); ok {
		r1 = rf(bill)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCostBreakdownByUser provides a mock function with given fields: bill
func (_m *BillStore) GetCostBreakdownByUser(bill *model.Bill) ([]model.UserCostBreakdown, error) {
	ret := _m.Called(bill)

	var r0 []model.UserCostBreakdown
	if rf, ok := ret.Get(0).(func(*model.Bill) []model.UserCostBreakdown); ok {
		r0 = rf(bill)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.UserCostBreakdown)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Bill) error); ok {
		r1 = rf(bill)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveFormal provides a mock function with given fields: bill, formalId
func (_m *BillStore) RemoveFormal(bill *model.Bill, formalId uuid.UUID) error {
	ret := _m.Called(bill, formalId)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Bill, uuid.UUID) error); ok {
		r0 = rf(bill, formalId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: bill
func (_m *BillStore) Update(bill *model.Bill) error {
	ret := _m.Called(bill)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Bill) error); ok {
		r0 = rf(bill)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewBillStore creates a new instance of BillStore. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewBillStore(t testing.TB) *BillStore {
	mock := &BillStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
