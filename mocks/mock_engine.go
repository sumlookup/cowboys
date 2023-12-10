// Package mock_engine is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockCowboysEngine is a mock of CowboysEngine interface.
type MockCowboysEngine struct {
	ctrl     *gomock.Controller
	recorder *MockCowboysEngineMockRecorder
}

// MockCowboysEngineMockRecorder is the mock recorder for MockCowboysEngine.
type MockCowboysEngineMockRecorder struct {
	mock *MockCowboysEngine
}

// NewMockCowboysEngine creates a new mock instance.
func NewMockCowboysEngine(ctrl *gomock.Controller) *MockCowboysEngine {
	mock := &MockCowboysEngine{ctrl: ctrl}
	mock.recorder = &MockCowboysEngineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCowboysEngine) EXPECT() *MockCowboysEngineMockRecorder {
	return m.recorder
}

// GameMode mocks base method.
func (m *MockCowboysEngine) GameMode() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GameMode")
	ret0, _ := ret[0].(string)
	return ret0
}

// GameMode indicates an expected call of GameMode.
func (mr *MockCowboysEngineMockRecorder) GameMode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GameMode", reflect.TypeOf((*MockCowboysEngine)(nil).GameMode))
}

// Run mocks base method.
func (m *MockCowboysEngine) Run(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockCowboysEngineMockRecorder) Run(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockCowboysEngine)(nil).Run), ctx)
}

// SetWinner mocks base method.
func (m *MockCowboysEngine) SetWinner(ctx context.Context, name string, winnerId, gameId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetWinner", ctx, name, winnerId, gameId)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetWinner indicates an expected call of SetWinner.
func (mr *MockCowboysEngineMockRecorder) SetWinner(ctx, name, winnerId, gameId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWinner", reflect.TypeOf((*MockCowboysEngine)(nil).SetWinner), ctx, name, winnerId, gameId)
}

// ShootRandomCowboy mocks base method.
func (m *MockCowboysEngine) ShootRandomCowboy(ctx context.Context, shooterId, gameId uuid.UUID, shooterName string, shooterDmg int32) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShootRandomCowboy", ctx, shooterId, gameId, shooterName, shooterDmg)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShootRandomCowboy indicates an expected call of ShootRandomCowboy.
func (mr *MockCowboysEngineMockRecorder) ShootRandomCowboy(ctx, shooterId, gameId, shooterName, shooterDmg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShootRandomCowboy", reflect.TypeOf((*MockCowboysEngine)(nil).ShootRandomCowboy), ctx, shooterId, gameId, shooterName, shooterDmg)
}

// StartGame mocks base method.
func (m *MockCowboysEngine) StartGame(ctx context.Context, id, gameId uuid.UUID, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartGame", ctx, id, gameId, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartGame indicates an expected call of StartGame.
func (mr *MockCowboysEngineMockRecorder) StartGame(ctx, id, gameId, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartGame", reflect.TypeOf((*MockCowboysEngine)(nil).StartGame), ctx, id, gameId, name)
}
