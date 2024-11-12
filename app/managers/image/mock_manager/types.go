// Code generated by MockGen. DO NOT EDIT.
// Source: types.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// IsURLAllowed mocks base method.
func (m *MockManager) IsURLAllowed(imageURL string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsURLAllowed", imageURL)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsURLAllowed indicates an expected call of IsURLAllowed.
func (mr *MockManagerMockRecorder) IsURLAllowed(imageURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsURLAllowed", reflect.TypeOf((*MockManager)(nil).IsURLAllowed), imageURL)
}

// ProcessImage mocks base method.
func (m *MockManager) ProcessImage(imageURL string, width, height int, format string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessImage", imageURL, width, height, format)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProcessImage indicates an expected call of ProcessImage.
func (mr *MockManagerMockRecorder) ProcessImage(imageURL, width, height, format interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessImage", reflect.TypeOf((*MockManager)(nil).ProcessImage), imageURL, width, height, format)
}