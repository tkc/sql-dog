// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/tkc/sql-dog/src/infrastructure/datastore/mysql (interfaces: GeneralLogRepository)

// Package mock_mysql is a generated GoMock package.
package mock_mysql

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockGeneralLogRepository is a mock of GeneralLogRepository interface
type MockGeneralLogRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGeneralLogRepositoryMockRecorder
}

// MockGeneralLogRepositoryMockRecorder is the mock recorder for MockGeneralLogRepository
type MockGeneralLogRepositoryMockRecorder struct {
	mock *MockGeneralLogRepository
}

// NewMockGeneralLogRepository creates a new mock instance
func NewMockGeneralLogRepository(ctrl *gomock.Controller) *MockGeneralLogRepository {
	mock := &MockGeneralLogRepository{ctrl: ctrl}
	mock.recorder = &MockGeneralLogRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGeneralLogRepository) EXPECT() *MockGeneralLogRepositoryMockRecorder {
	return m.recorder
}

// Clear mocks base method
func (m *MockGeneralLogRepository) Clear() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clear")
	ret0, _ := ret[0].(error)
	return ret0
}

// Clear indicates an expected call of Clear
func (mr *MockGeneralLogRepositoryMockRecorder) Clear() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockGeneralLogRepository)(nil).Clear))
}

// GetQueries mocks base method
func (m *MockGeneralLogRepository) GetQueries() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQueries")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQueries indicates an expected call of GetQueries
func (mr *MockGeneralLogRepositoryMockRecorder) GetQueries() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQueries", reflect.TypeOf((*MockGeneralLogRepository)(nil).GetQueries))
}
