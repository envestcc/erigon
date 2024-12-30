// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ledgerwatch/erigon/v2/polygon/sync (interfaces: Storage)
//
// Generated by this command:
//
//	mockgen -typed=true -destination=./storage_mock.go -package=sync . Storage
//

// Package sync is a generated GoMock package.
package sync

import (
	context "context"
	reflect "reflect"

	types "github.com/ledgerwatch/erigon/v2/core/types"
	gomock "go.uber.org/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Flush mocks base method.
func (m *MockStorage) Flush(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Flush", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Flush indicates an expected call of Flush.
func (mr *MockStorageMockRecorder) Flush(arg0 any) *MockStorageFlushCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flush", reflect.TypeOf((*MockStorage)(nil).Flush), arg0)
	return &MockStorageFlushCall{Call: call}
}

// MockStorageFlushCall wrap *gomock.Call
type MockStorageFlushCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStorageFlushCall) Return(arg0 error) *MockStorageFlushCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStorageFlushCall) Do(f func(context.Context) error) *MockStorageFlushCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStorageFlushCall) DoAndReturn(f func(context.Context) error) *MockStorageFlushCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// InsertBlocks mocks base method.
func (m *MockStorage) InsertBlocks(arg0 context.Context, arg1 []*types.Block) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertBlocks", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertBlocks indicates an expected call of InsertBlocks.
func (mr *MockStorageMockRecorder) InsertBlocks(arg0, arg1 any) *MockStorageInsertBlocksCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertBlocks", reflect.TypeOf((*MockStorage)(nil).InsertBlocks), arg0, arg1)
	return &MockStorageInsertBlocksCall{Call: call}
}

// MockStorageInsertBlocksCall wrap *gomock.Call
type MockStorageInsertBlocksCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStorageInsertBlocksCall) Return(arg0 error) *MockStorageInsertBlocksCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStorageInsertBlocksCall) Do(f func(context.Context, []*types.Block) error) *MockStorageInsertBlocksCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStorageInsertBlocksCall) DoAndReturn(f func(context.Context, []*types.Block) error) *MockStorageInsertBlocksCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Run mocks base method.
func (m *MockStorage) Run(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockStorageMockRecorder) Run(arg0 any) *MockStorageRunCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockStorage)(nil).Run), arg0)
	return &MockStorageRunCall{Call: call}
}

// MockStorageRunCall wrap *gomock.Call
type MockStorageRunCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStorageRunCall) Return(arg0 error) *MockStorageRunCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStorageRunCall) Do(f func(context.Context) error) *MockStorageRunCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStorageRunCall) DoAndReturn(f func(context.Context) error) *MockStorageRunCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
