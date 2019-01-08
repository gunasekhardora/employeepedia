package main

import (
	"github.com/stretchr/testify/mock"
)

// MockStore contains additonal methods for inspection
type MockStore struct {
	mock.Mock
}

/*
	When this method is called, `m.Called` records the call, and also
	returns the result that we pass to it (which you will see in the
	handler tests)
*/
func (m *MockStore) CreateEmployee(employee *Employee) error {

	rets := m.Called(employee)
	return rets.Error(0)
}

/*
	Since `rets.Get()` is a generic method, that returns whatever we pass to it,
	we need to typecast it to the type we expect, which in this case is []*Employee
*/
func (m *MockStore) GetEmployees() ([]*Employee, error) {
	rets := m.Called()

	return rets.Get(0).([]*Employee), rets.Error(1)
}

// InitMockStore function assigns
// a new MockStore instance to it, instead of an actual store
func InitMockStore() *MockStore {

	s := new(MockStore)
	store = s
	return s
}
