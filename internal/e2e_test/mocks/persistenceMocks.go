// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/persistence/persistence.go

// Package mocks is a generated GoMock package.
package mocks

import (
	persistence "github.com/florianloch/cassette/internal/persistence"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPlayerStatesPersistor is a mock of PlayerStatesPersistor interface
type MockPlayerStatesPersistor struct {
	ctrl     *gomock.Controller
	recorder *MockPlayerStatesPersistorMockRecorder
}

// MockPlayerStatesPersistorMockRecorder is the mock recorder for MockPlayerStatesPersistor
type MockPlayerStatesPersistorMockRecorder struct {
	mock *MockPlayerStatesPersistor
}

// NewMockPlayerStatesPersistor creates a new mock instance
func NewMockPlayerStatesPersistor(ctrl *gomock.Controller) *MockPlayerStatesPersistor {
	mock := &MockPlayerStatesPersistor{ctrl: ctrl}
	mock.recorder = &MockPlayerStatesPersistorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPlayerStatesPersistor) EXPECT() *MockPlayerStatesPersistorMockRecorder {
	return m.recorder
}

// LoadPlayerStates mocks base method
func (m *MockPlayerStatesPersistor) LoadPlayerStates(userID string) ([]*persistence.PlayerState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadPlayerStates", userID)
	ret0, _ := ret[0].([]*persistence.PlayerState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadPlayerStates indicates an expected call of LoadPlayerStates
func (mr *MockPlayerStatesPersistorMockRecorder) LoadPlayerStates(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadPlayerStates", reflect.TypeOf((*MockPlayerStatesPersistor)(nil).LoadPlayerStates), userID)
}

// SavePlayerStates mocks base method
func (m *MockPlayerStatesPersistor) SavePlayerStates(userID string, playerStates []*persistence.PlayerState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SavePlayerStates", userID, playerStates)
	ret0, _ := ret[0].(error)
	return ret0
}

// SavePlayerStates indicates an expected call of SavePlayerStates
func (mr *MockPlayerStatesPersistorMockRecorder) SavePlayerStates(userID, playerStates interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SavePlayerStates", reflect.TypeOf((*MockPlayerStatesPersistor)(nil).SavePlayerStates), userID, playerStates)
}

// FetchJSONDump mocks base method
func (m *MockPlayerStatesPersistor) FetchJSONDump(userID string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchJSONDump", userID)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchJSONDump indicates an expected call of FetchJSONDump
func (mr *MockPlayerStatesPersistorMockRecorder) FetchJSONDump(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchJSONDump", reflect.TypeOf((*MockPlayerStatesPersistor)(nil).FetchJSONDump), userID)
}

// DeleteUserRecord mocks base method
func (m *MockPlayerStatesPersistor) DeleteUserRecord(userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserRecord", userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserRecord indicates an expected call of DeleteUserRecord
func (mr *MockPlayerStatesPersistorMockRecorder) DeleteUserRecord(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserRecord", reflect.TypeOf((*MockPlayerStatesPersistor)(nil).DeleteUserRecord), userID)
}
