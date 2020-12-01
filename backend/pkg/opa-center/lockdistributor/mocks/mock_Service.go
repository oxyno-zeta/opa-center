// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/oxyno-zeta/opa-center/pkg/opa-center/lockdistributor (interfaces: Service)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	lockdistributor "github.com/oxyno-zeta/opa-center/pkg/opa-center/lockdistributor"
	log "github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetLock mocks base method
func (m *MockService) GetLock(arg0 string) lockdistributor.Lock {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLock", arg0)
	ret0, _ := ret[0].(lockdistributor.Lock)
	return ret0
}

// GetLock indicates an expected call of GetLock
func (mr *MockServiceMockRecorder) GetLock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLock", reflect.TypeOf((*MockService)(nil).GetLock), arg0)
}

// Initialize mocks base method
func (m *MockService) Initialize(arg0 log.Logger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize
func (mr *MockServiceMockRecorder) Initialize(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockService)(nil).Initialize), arg0)
}
