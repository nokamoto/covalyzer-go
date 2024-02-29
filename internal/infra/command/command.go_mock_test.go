// Code generated by MockGen. DO NOT EDIT.
// Source: command.go
//
// Generated by this command:
//
//	mockgen -source=command.go -destination=command.go_mock_test.go -package=command
//

// Package command is a generated GoMock package.
package command

import (
	bytes "bytes"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// Mockrunner is a mock of runner interface.
type Mockrunner struct {
	ctrl     *gomock.Controller
	recorder *MockrunnerMockRecorder
}

// MockrunnerMockRecorder is the mock recorder for Mockrunner.
type MockrunnerMockRecorder struct {
	mock *Mockrunner
}

// NewMockrunner creates a new mock instance.
func NewMockrunner(ctrl *gomock.Controller) *Mockrunner {
	mock := &Mockrunner{ctrl: ctrl}
	mock.recorder = &MockrunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrunner) EXPECT() *MockrunnerMockRecorder {
	return m.recorder
}

// run mocks base method.
func (m *Mockrunner) run(cmd string, args []string, opts ...option) error {
	m.ctrl.T.Helper()
	varargs := []any{cmd, args}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "run", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// run indicates an expected call of run.
func (mr *MockrunnerMockRecorder) run(cmd, args any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{cmd, args}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "run", reflect.TypeOf((*Mockrunner)(nil).run), varargs...)
}

// runO mocks base method.
func (m *Mockrunner) runO(cmd string, args []string, opts ...option) (*bytes.Buffer, error) {
	m.ctrl.T.Helper()
	varargs := []any{cmd, args}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "runO", varargs...)
	ret0, _ := ret[0].(*bytes.Buffer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// runO indicates an expected call of runO.
func (mr *MockrunnerMockRecorder) runO(cmd, args any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{cmd, args}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "runO", reflect.TypeOf((*Mockrunner)(nil).runO), varargs...)
}