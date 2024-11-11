// Code generated by MockGen. DO NOT EDIT.
// Source: types.go

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// HandleIndex mocks base method.
func (m *MockHandler) HandleIndex(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleIndex", c)
}

// HandleIndex indicates an expected call of HandleIndex.
func (mr *MockHandlerMockRecorder) HandleIndex(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleIndex", reflect.TypeOf((*MockHandler)(nil).HandleIndex), c)
}
