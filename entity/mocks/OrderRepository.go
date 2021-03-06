// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	entity "github.com/hieronimusbudi/komodo-backend/entity"
	mock "github.com/stretchr/testify/mock"

	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
)

// OrderRepository is an autogenerated mock type for the OrderRepository type
type OrderRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: order
func (_m *OrderRepository) Delete(order *entity.Order) resterrors.RestErr {
	ret := _m.Called(order)

	var r0 resterrors.RestErr
	if rf, ok := ret.Get(0).(func(*entity.Order) resterrors.RestErr); ok {
		r0 = rf(order)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(resterrors.RestErr)
		}
	}

	return r0
}

// GetAll provides a mock function with given fields:
func (_m *OrderRepository) GetAll() ([]entity.Order, resterrors.RestErr) {
	ret := _m.Called()

	var r0 []entity.Order
	if rf, ok := ret.Get(0).(func() []entity.Order); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Order)
		}
	}

	var r1 resterrors.RestErr
	if rf, ok := ret.Get(1).(func() resterrors.RestErr); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(resterrors.RestErr)
		}
	}

	return r0, r1
}

// GetByBuyerID provides a mock function with given fields: buyerID
func (_m *OrderRepository) GetByBuyerID(buyerID int64) ([]entity.Order, resterrors.RestErr) {
	ret := _m.Called(buyerID)

	var r0 []entity.Order
	if rf, ok := ret.Get(0).(func(int64) []entity.Order); ok {
		r0 = rf(buyerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Order)
		}
	}

	var r1 resterrors.RestErr
	if rf, ok := ret.Get(1).(func(int64) resterrors.RestErr); ok {
		r1 = rf(buyerID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(resterrors.RestErr)
		}
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: order
func (_m *OrderRepository) GetByID(order *entity.Order) (entity.Order, resterrors.RestErr) {
	ret := _m.Called(order)

	var r0 entity.Order
	if rf, ok := ret.Get(0).(func(*entity.Order) entity.Order); ok {
		r0 = rf(order)
	} else {
		r0 = ret.Get(0).(entity.Order)
	}

	var r1 resterrors.RestErr
	if rf, ok := ret.Get(1).(func(*entity.Order) resterrors.RestErr); ok {
		r1 = rf(order)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(resterrors.RestErr)
		}
	}

	return r0, r1
}

// GetBySellerID provides a mock function with given fields: buyerID
func (_m *OrderRepository) GetBySellerID(buyerID int64) ([]entity.Order, resterrors.RestErr) {
	ret := _m.Called(buyerID)

	var r0 []entity.Order
	if rf, ok := ret.Get(0).(func(int64) []entity.Order); ok {
		r0 = rf(buyerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Order)
		}
	}

	var r1 resterrors.RestErr
	if rf, ok := ret.Get(1).(func(int64) resterrors.RestErr); ok {
		r1 = rf(buyerID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(resterrors.RestErr)
		}
	}

	return r0, r1
}

// Store provides a mock function with given fields: order
func (_m *OrderRepository) Store(order *entity.Order) resterrors.RestErr {
	ret := _m.Called(order)

	var r0 resterrors.RestErr
	if rf, ok := ret.Get(0).(func(*entity.Order) resterrors.RestErr); ok {
		r0 = rf(order)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(resterrors.RestErr)
		}
	}

	return r0
}

// Update provides a mock function with given fields: order
func (_m *OrderRepository) Update(order *entity.Order) resterrors.RestErr {
	ret := _m.Called(order)

	var r0 resterrors.RestErr
	if rf, ok := ret.Get(0).(func(*entity.Order) resterrors.RestErr); ok {
		r0 = rf(order)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(resterrors.RestErr)
		}
	}

	return r0
}
