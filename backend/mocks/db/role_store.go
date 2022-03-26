// Code generated by mockery 2.10.0. DO NOT EDIT.

package mocks

import (
	model "github.com/kcsu/store/model"
	mock "github.com/stretchr/testify/mock"
)

// RoleStore is an autogenerated mock type for the RoleStore type
type RoleStore struct {
	mock.Mock
}

// AddUserRole provides a mock function with given fields: role, user
func (_m *RoleStore) AddUserRole(role *model.Role, user *model.User) error {
	ret := _m.Called(role, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Role, *model.User) error); ok {
		r0 = rf(role, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: role
func (_m *RoleStore) Create(role *model.Role) error {
	ret := _m.Called(role)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Role) error); ok {
		r0 = rf(role)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreatePermission provides a mock function with given fields: permission
func (_m *RoleStore) CreatePermission(permission *model.Permission) error {
	ret := _m.Called(permission)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Permission) error); ok {
		r0 = rf(permission)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePermission provides a mock function with given fields: id
func (_m *RoleStore) DeletePermission(id int) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: id
func (_m *RoleStore) Find(id int) (model.Role, error) {
	ret := _m.Called(id)

	var r0 model.Role
	if rf, ok := ret.Get(0).(func(int) model.Role); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.Role)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields:
func (_m *RoleStore) Get() ([]model.Role, error) {
	ret := _m.Called()

	var r0 []model.Role
	if rf, ok := ret.Get(0).(func() []model.Role); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Role)
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

// GetUserRoles provides a mock function with given fields:
func (_m *RoleStore) GetUserRoles() ([]model.UserRole, error) {
	ret := _m.Called()

	var r0 []model.UserRole
	if rf, ok := ret.Get(0).(func() []model.UserRole); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.UserRole)
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

// RemoveUserRole provides a mock function with given fields: role, user
func (_m *RoleStore) RemoveUserRole(role *model.Role, user *model.User) error {
	ret := _m.Called(role, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Role, *model.User) error); ok {
		r0 = rf(role, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
