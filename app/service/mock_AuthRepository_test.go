// Code generated by mockery v2.43.2. DO NOT EDIT.

package service

import (
	context "context"
	entity "go-dating-app/app/entity"

	mock "github.com/stretchr/testify/mock"
)

// MockAuthRepository is an autogenerated mock type for the AuthRepository type
type MockAuthRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, user
func (_m *MockAuthRepository) Create(ctx context.Context, user *entity.User) error {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByEmail provides a mock function with given fields: ctx, email
func (_m *MockAuthRepository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for FindByEmail")
	}

	var r0 entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (entity.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) entity.User); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(entity.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockAuthRepository creates a new instance of MockAuthRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAuthRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAuthRepository {
	mock := &MockAuthRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
