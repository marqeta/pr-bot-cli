// Code generated by mockery v2.33.0. DO NOT EDIT.

package evaluation

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockManager is an autogenerated mock type for the Manager type
type MockManager struct {
	mock.Mock
}

type MockManager_Expecter struct {
	mock *mock.Mock
}

func (_m *MockManager) EXPECT() *MockManager_Expecter {
	return &MockManager_Expecter{mock: &_m.Mock}
}

// GetReport provides a mock function with given fields: ctx, pr, deliveryID
func (_m *MockManager) GetReport(ctx context.Context, pr string, deliveryID string) (*Report, error) {
	ret := _m.Called(ctx, pr, deliveryID)

	var r0 *Report
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*Report, error)); ok {
		return rf(ctx, pr, deliveryID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *Report); ok {
		r0 = rf(ctx, pr, deliveryID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Report)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, pr, deliveryID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockManager_GetReport_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetReport'
type MockManager_GetReport_Call struct {
	*mock.Call
}

// GetReport is a helper method to define mock.On call
//   - ctx context.Context
//   - pr string
//   - deliveryID string
func (_e *MockManager_Expecter) GetReport(ctx interface{}, pr interface{}, deliveryID interface{}) *MockManager_GetReport_Call {
	return &MockManager_GetReport_Call{Call: _e.mock.On("GetReport", ctx, pr, deliveryID)}
}

func (_c *MockManager_GetReport_Call) Run(run func(ctx context.Context, pr string, deliveryID string)) *MockManager_GetReport_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockManager_GetReport_Call) Return(_a0 *Report, _a1 error) *MockManager_GetReport_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockManager_GetReport_Call) RunAndReturn(run func(context.Context, string, string) (*Report, error)) *MockManager_GetReport_Call {
	_c.Call.Return(run)
	return _c
}

// ListReports provides a mock function with given fields: ctx, pr
func (_m *MockManager) ListReports(ctx context.Context, pr string) ([]ReportMetadata, error) {
	ret := _m.Called(ctx, pr)

	var r0 []ReportMetadata
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]ReportMetadata, error)); ok {
		return rf(ctx, pr)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []ReportMetadata); ok {
		r0 = rf(ctx, pr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ReportMetadata)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, pr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockManager_ListReports_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListReports'
type MockManager_ListReports_Call struct {
	*mock.Call
}

// ListReports is a helper method to define mock.On call
//   - ctx context.Context
//   - pr string
func (_e *MockManager_Expecter) ListReports(ctx interface{}, pr interface{}) *MockManager_ListReports_Call {
	return &MockManager_ListReports_Call{Call: _e.mock.On("ListReports", ctx, pr)}
}

func (_c *MockManager_ListReports_Call) Run(run func(ctx context.Context, pr string)) *MockManager_ListReports_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockManager_ListReports_Call) Return(_a0 []ReportMetadata, _a1 error) *MockManager_ListReports_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockManager_ListReports_Call) RunAndReturn(run func(context.Context, string) ([]ReportMetadata, error)) *MockManager_ListReports_Call {
	_c.Call.Return(run)
	return _c
}

// NewReportBuilder provides a mock function with given fields: ctx, pr, reqID, deliveryID
func (_m *MockManager) NewReportBuilder(ctx context.Context, pr string, reqID string, deliveryID string) ReportBuilder {
	ret := _m.Called(ctx, pr, reqID, deliveryID)

	var r0 ReportBuilder
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) ReportBuilder); ok {
		r0 = rf(ctx, pr, reqID, deliveryID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ReportBuilder)
		}
	}

	return r0
}

// MockManager_NewReportBuilder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewReportBuilder'
type MockManager_NewReportBuilder_Call struct {
	*mock.Call
}

// NewReportBuilder is a helper method to define mock.On call
//   - ctx context.Context
//   - pr string
//   - reqID string
//   - deliveryID string
func (_e *MockManager_Expecter) NewReportBuilder(ctx interface{}, pr interface{}, reqID interface{}, deliveryID interface{}) *MockManager_NewReportBuilder_Call {
	return &MockManager_NewReportBuilder_Call{Call: _e.mock.On("NewReportBuilder", ctx, pr, reqID, deliveryID)}
}

func (_c *MockManager_NewReportBuilder_Call) Run(run func(ctx context.Context, pr string, reqID string, deliveryID string)) *MockManager_NewReportBuilder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockManager_NewReportBuilder_Call) Return(_a0 ReportBuilder) *MockManager_NewReportBuilder_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_NewReportBuilder_Call) RunAndReturn(run func(context.Context, string, string, string) ReportBuilder) *MockManager_NewReportBuilder_Call {
	_c.Call.Return(run)
	return _c
}

// StoreReport provides a mock function with given fields: ctx, builder
func (_m *MockManager) StoreReport(ctx context.Context, builder ReportBuilder) error {
	ret := _m.Called(ctx, builder)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ReportBuilder) error); ok {
		r0 = rf(ctx, builder)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockManager_StoreReport_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StoreReport'
type MockManager_StoreReport_Call struct {
	*mock.Call
}

// StoreReport is a helper method to define mock.On call
//   - ctx context.Context
//   - builder ReportBuilder
func (_e *MockManager_Expecter) StoreReport(ctx interface{}, builder interface{}) *MockManager_StoreReport_Call {
	return &MockManager_StoreReport_Call{Call: _e.mock.On("StoreReport", ctx, builder)}
}

func (_c *MockManager_StoreReport_Call) Run(run func(ctx context.Context, builder ReportBuilder)) *MockManager_StoreReport_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ReportBuilder))
	})
	return _c
}

func (_c *MockManager_StoreReport_Call) Return(_a0 error) *MockManager_StoreReport_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_StoreReport_Call) RunAndReturn(run func(context.Context, ReportBuilder) error) *MockManager_StoreReport_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockManager creates a new instance of MockManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockManager {
	mock := &MockManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}