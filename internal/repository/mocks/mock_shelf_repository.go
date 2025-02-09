// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Users\NetBoock\Desktop\bookshelf\internal\repository\shelf_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "myproject/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockShelfRepository is a mock of ShelfRepository interface.
type MockShelfRepository struct {
	ctrl     *gomock.Controller
	recorder *MockShelfRepositoryMockRecorder
}

// MockShelfRepositoryMockRecorder is the mock recorder for MockShelfRepository.
type MockShelfRepositoryMockRecorder struct {
	mock *MockShelfRepository
}

// NewMockShelfRepository creates a new mock instance.
func NewMockShelfRepository(ctrl *gomock.Controller) *MockShelfRepository {
	mock := &MockShelfRepository{ctrl: ctrl}
	mock.recorder = &MockShelfRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShelfRepository) EXPECT() *MockShelfRepositoryMockRecorder {
	return m.recorder
}

// AddShelf mocks base method.
func (m *MockShelfRepository) AddShelf() (*models.Shelf, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddShelf")
	ret0, _ := ret[0].(*models.Shelf)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddShelf indicates an expected call of AddShelf.
func (mr *MockShelfRepositoryMockRecorder) AddShelf() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddShelf", reflect.TypeOf((*MockShelfRepository)(nil).AddShelf))
}

// DeleteShelf mocks base method.
func (m *MockShelfRepository) DeleteShelf(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteShelf", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteShelf indicates an expected call of DeleteShelf.
func (mr *MockShelfRepositoryMockRecorder) DeleteShelf(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteShelf", reflect.TypeOf((*MockShelfRepository)(nil).DeleteShelf), id)
}

// GetShelf mocks base method.
func (m *MockShelfRepository) GetShelf(id int) (*models.Shelf, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShelf", id)
	ret0, _ := ret[0].(*models.Shelf)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShelf indicates an expected call of GetShelf.
func (mr *MockShelfRepositoryMockRecorder) GetShelf(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShelf", reflect.TypeOf((*MockShelfRepository)(nil).GetShelf), id)
}
