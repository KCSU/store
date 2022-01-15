// Code generated by mockery 2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	idtoken "google.golang.org/api/idtoken"
)

// IdTokenValidator is an autogenerated mock type for the IdTokenValidator type
type IdTokenValidator struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0, _a1, _a2
func (_m *IdTokenValidator) Execute(_a0 context.Context, _a1 string, _a2 string) (*idtoken.Payload, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *idtoken.Payload
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *idtoken.Payload); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*idtoken.Payload)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
