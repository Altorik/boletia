// Code generated by mockery v2.42.0. DO NOT EDIT.

package storagemocks

import (
	bole "boletia/api/internal"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IDatabaseRepository is an autogenerated mock type for the IDatabaseRepository type
type IDatabaseRepository struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, criteria
func (_m *IDatabaseRepository) Get(ctx context.Context, criteria bole.Criteria) ([]bole.Currency, error) {
	ret := _m.Called(ctx, criteria)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 []bole.Currency
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, bole.Criteria) ([]bole.Currency, error)); ok {
		return rf(ctx, criteria)
	}
	if rf, ok := ret.Get(0).(func(context.Context, bole.Criteria) []bole.Currency); ok {
		r0 = rf(ctx, criteria)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]bole.Currency)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, bole.Criteria) error); ok {
		r1 = rf(ctx, criteria)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIDatabaseRepository creates a new instance of IDatabaseRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIDatabaseRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IDatabaseRepository {
	mock := &IDatabaseRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}