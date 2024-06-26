// Code generated by MockGen. DO NOT EDIT.
// Source: internal/users/usecase.go
//
// Generated by this command:
//
//	mockgen -source internal/users/usecase.go -destination internal/users/mock/usecase_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	dtos "github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
	gomock "go.uber.org/mock/gomock"
)

// MockUsecase is a mock of Usecase interface.
type MockUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUsecaseMockRecorder
}

// MockUsecaseMockRecorder is the mock recorder for MockUsecase.
type MockUsecaseMockRecorder struct {
	mock *MockUsecase
}

// NewMockUsecase creates a new mock instance.
func NewMockUsecase(ctrl *gomock.Controller) *MockUsecase {
	mock := &MockUsecase{ctrl: ctrl}
	mock.recorder = &MockUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsecase) EXPECT() *MockUsecaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUsecase) Create(ctx context.Context, payload dtos.CreateUserRequest) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, payload)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUsecaseMockRecorder) Create(ctx, payload any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUsecase)(nil).Create), ctx, payload)
}

// Detail mocks base method.
func (m *MockUsecase) Detail(ctx context.Context, id int64) (dtos.UserDetailResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Detail", ctx, id)
	ret0, _ := ret[0].(dtos.UserDetailResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Detail indicates an expected call of Detail.
func (mr *MockUsecaseMockRecorder) Detail(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Detail", reflect.TypeOf((*MockUsecase)(nil).Detail), ctx, id)
}

// Login mocks base method.
func (m *MockUsecase) Login(ctx context.Context, request dtos.UserLoginRequest) (dtos.UserLoginResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, request)
	ret0, _ := ret[0].(dtos.UserLoginResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUsecaseMockRecorder) Login(ctx, request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUsecase)(nil).Login), ctx, request)
}

// PartialUpdate mocks base method.
func (m *MockUsecase) PartialUpdate(ctx context.Context, data dtos.UpdateUserRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PartialUpdate", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// PartialUpdate indicates an expected call of PartialUpdate.
func (mr *MockUsecaseMockRecorder) PartialUpdate(ctx, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PartialUpdate", reflect.TypeOf((*MockUsecase)(nil).PartialUpdate), ctx, data)
}

// UpdateStatus mocks base method.
func (m *MockUsecase) UpdateStatus(ctx context.Context, req dtos.UpdateUserStatusRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockUsecaseMockRecorder) UpdateStatus(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockUsecase)(nil).UpdateStatus), ctx, req)
}
