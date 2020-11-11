// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authentication (interfaces: Client)

// Package mocks is a generated GoMock package.
package mocks

import (
	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	regexp "regexp"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Middleware mocks base method
func (m *MockClient) Middleware(arg0 []*regexp.Regexp) gin.HandlerFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Middleware", arg0)
	ret0, _ := ret[0].(gin.HandlerFunc)
	return ret0
}

// Middleware indicates an expected call of Middleware
func (mr *MockClientMockRecorder) Middleware(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Middleware", reflect.TypeOf((*MockClient)(nil).Middleware), arg0)
}

// OIDCEndpoints mocks base method
func (m *MockClient) OIDCEndpoints(arg0 gin.IRouter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OIDCEndpoints", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// OIDCEndpoints indicates an expected call of OIDCEndpoints
func (mr *MockClientMockRecorder) OIDCEndpoints(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OIDCEndpoints", reflect.TypeOf((*MockClient)(nil).OIDCEndpoints), arg0)
}
