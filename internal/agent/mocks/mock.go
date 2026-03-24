// Package mocks
package mocks

import (
	"reflect"

	"github.com/golang/mock/gomock"
)

// MockMemStorage is a mock of MemStorage interface
type MockMemStorage struct {
	ctrl     *gomock.Controller
	recorder *MockMemStorageMockRecorder
}

// MockMemStorageMockRecorder is the mock recorder for MockMemStorage
type MockMemStorageMockRecorder struct {
	mock *MockMemStorage
}

// NewMockMemStorage creates a new mock instance
func NewMockMemStorage(ctrl *gomock.Controller) *MockMemStorage {
	mock := &MockMemStorage{ctrl: ctrl}
	mock.recorder = &MockMemStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMemStorage) EXPECT() *MockMemStorageMockRecorder {
	return m.recorder
}

// SetGauge mocks base method
func (m *MockMemStorage) SetGauge(key string, value float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetGauge", key, value)
}

// SetGauge indicates an expected call of SetGauge
func (mr *MockMemStorageMockRecorder) SetGauge(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetGauge", reflect.TypeOf((*MockMemStorage)(nil).SetGauge), key, value)
}

// SetCounter mocks base method
func (m *MockMemStorage) SetCounter(key string, value int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetCounter", key, value)
}

// SetCounter indicates an expected call of SetCounter
func (mr *MockMemStorageMockRecorder) SetCounter(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCounter", reflect.TypeOf((*MockMemStorage)(nil).SetCounter), key, value)
}

// GetGauge mocks base method
func (m *MockMemStorage) GetGauge() map[string]float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGauge")
	ret0, _ := ret[0].(map[string]float64)
	return ret0
}

// GetGauge indicates an expected call of GetGauge
func (mr *MockMemStorageMockRecorder) GetGauge() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGauge", reflect.TypeOf((*MockMemStorage)(nil).GetGauge))
}

// GetCounter mocks base method
func (m *MockMemStorage) GetCounter() map[string]int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounter")
	ret0, _ := ret[0].(map[string]int64)
	return ret0
}

// GetCounter indicates an expected call of GetCounter
func (mr *MockMemStorageMockRecorder) GetCounter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounter", reflect.TypeOf((*MockMemStorage)(nil).GetCounter))
}

// AddCounter mocks base method
func (m *MockMemStorage) AddCounter(key string, value int64) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCounter", key, value)
	ret0, _ := ret[0].(bool)
	return ret0
}

// AddCounter indicates an expected call of AddCounter
func (mr *MockMemStorageMockRecorder) AddCounter(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCounter", reflect.TypeOf((*MockMemStorage)(nil).AddCounter), key, value)
}
