// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ledgerwatch/erigon/v2/cl/phase1/network/services (interfaces: BLSToExecutionChangeService)
//
// Generated by this command:
//
//	mockgen -typed=true -destination=./mock_services/bls_to_execution_change_service_mock.go -package=mock_services . BLSToExecutionChangeService
//

// Package mock_services is a generated GoMock package.
package mock_services

import (
	context "context"
	reflect "reflect"

	cltypes "github.com/ledgerwatch/erigon/v2/cl/cltypes"
	gomock "go.uber.org/mock/gomock"
)

// MockBLSToExecutionChangeService is a mock of BLSToExecutionChangeService interface.
type MockBLSToExecutionChangeService struct {
	ctrl     *gomock.Controller
	recorder *MockBLSToExecutionChangeServiceMockRecorder
}

// MockBLSToExecutionChangeServiceMockRecorder is the mock recorder for MockBLSToExecutionChangeService.
type MockBLSToExecutionChangeServiceMockRecorder struct {
	mock *MockBLSToExecutionChangeService
}

// NewMockBLSToExecutionChangeService creates a new mock instance.
func NewMockBLSToExecutionChangeService(ctrl *gomock.Controller) *MockBLSToExecutionChangeService {
	mock := &MockBLSToExecutionChangeService{ctrl: ctrl}
	mock.recorder = &MockBLSToExecutionChangeServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBLSToExecutionChangeService) EXPECT() *MockBLSToExecutionChangeServiceMockRecorder {
	return m.recorder
}

// ProcessMessage mocks base method.
func (m *MockBLSToExecutionChangeService) ProcessMessage(arg0 context.Context, arg1 *uint64, arg2 *cltypes.SignedBLSToExecutionChange) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessMessage", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessMessage indicates an expected call of ProcessMessage.
func (mr *MockBLSToExecutionChangeServiceMockRecorder) ProcessMessage(arg0, arg1, arg2 any) *MockBLSToExecutionChangeServiceProcessMessageCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessMessage", reflect.TypeOf((*MockBLSToExecutionChangeService)(nil).ProcessMessage), arg0, arg1, arg2)
	return &MockBLSToExecutionChangeServiceProcessMessageCall{Call: call}
}

// MockBLSToExecutionChangeServiceProcessMessageCall wrap *gomock.Call
type MockBLSToExecutionChangeServiceProcessMessageCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockBLSToExecutionChangeServiceProcessMessageCall) Return(arg0 error) *MockBLSToExecutionChangeServiceProcessMessageCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockBLSToExecutionChangeServiceProcessMessageCall) Do(f func(context.Context, *uint64, *cltypes.SignedBLSToExecutionChange) error) *MockBLSToExecutionChangeServiceProcessMessageCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockBLSToExecutionChangeServiceProcessMessageCall) DoAndReturn(f func(context.Context, *uint64, *cltypes.SignedBLSToExecutionChange) error) *MockBLSToExecutionChangeServiceProcessMessageCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
