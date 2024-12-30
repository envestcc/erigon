// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ledgerwatch/erigon/v2/polygon/heimdall (interfaces: HttpClient)
//
// Generated by this command:
//
//	mockgen -typed=true -destination=./http_client_mock.go -package=heimdall . HttpClient
//

// Package heimdall is a generated GoMock package.
package heimdall

import (
	http "net/http"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockHttpClient is a mock of HttpClient interface.
type MockHttpClient struct {
	ctrl     *gomock.Controller
	recorder *MockHttpClientMockRecorder
}

// MockHttpClientMockRecorder is the mock recorder for MockHttpClient.
type MockHttpClientMockRecorder struct {
	mock *MockHttpClient
}

// NewMockHttpClient creates a new mock instance.
func NewMockHttpClient(ctrl *gomock.Controller) *MockHttpClient {
	mock := &MockHttpClient{ctrl: ctrl}
	mock.recorder = &MockHttpClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHttpClient) EXPECT() *MockHttpClientMockRecorder {
	return m.recorder
}

// CloseIdleConnections mocks base method.
func (m *MockHttpClient) CloseIdleConnections() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CloseIdleConnections")
}

// CloseIdleConnections indicates an expected call of CloseIdleConnections.
func (mr *MockHttpClientMockRecorder) CloseIdleConnections() *MockHttpClientCloseIdleConnectionsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseIdleConnections", reflect.TypeOf((*MockHttpClient)(nil).CloseIdleConnections))
	return &MockHttpClientCloseIdleConnectionsCall{Call: call}
}

// MockHttpClientCloseIdleConnectionsCall wrap *gomock.Call
type MockHttpClientCloseIdleConnectionsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockHttpClientCloseIdleConnectionsCall) Return() *MockHttpClientCloseIdleConnectionsCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockHttpClientCloseIdleConnectionsCall) Do(f func()) *MockHttpClientCloseIdleConnectionsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockHttpClientCloseIdleConnectionsCall) DoAndReturn(f func()) *MockHttpClientCloseIdleConnectionsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Do mocks base method.
func (m *MockHttpClient) Do(arg0 *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockHttpClientMockRecorder) Do(arg0 any) *MockHttpClientDoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockHttpClient)(nil).Do), arg0)
	return &MockHttpClientDoCall{Call: call}
}

// MockHttpClientDoCall wrap *gomock.Call
type MockHttpClientDoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockHttpClientDoCall) Return(arg0 *http.Response, arg1 error) *MockHttpClientDoCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockHttpClientDoCall) Do(f func(*http.Request) (*http.Response, error)) *MockHttpClientDoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockHttpClientDoCall) DoAndReturn(f func(*http.Request) (*http.Response, error)) *MockHttpClientDoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
