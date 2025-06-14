// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go
//
// Generated by this command:
//
//	mockgen -source=interface.go -destination=mocks/mock-slack.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	slack "github.com/pal-paul/go-libraries/pkg/slack"
	gomock "go.uber.org/mock/gomock"
)

// MockISlack is a mock of ISlack interface.
type MockISlack struct {
	ctrl     *gomock.Controller
	recorder *MockISlackMockRecorder
	isgomock struct{}
}

// MockISlackMockRecorder is the mock recorder for MockISlack.
type MockISlackMockRecorder struct {
	mock *MockISlack
}

// NewMockISlack creates a new mock instance.
func NewMockISlack(ctrl *gomock.Controller) *MockISlack {
	mock := &MockISlack{ctrl: ctrl}
	mock.recorder = &MockISlackMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockISlack) EXPECT() *MockISlackMockRecorder {
	return m.recorder
}

// AddFormattedMessage mocks base method.
func (m *MockISlack) AddFormattedMessage(channel string, message slack.Message) (slack.MessageRef, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFormattedMessage", channel, message)
	ret0, _ := ret[0].(slack.MessageRef)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddFormattedMessage indicates an expected call of AddFormattedMessage.
func (mr *MockISlackMockRecorder) AddFormattedMessage(channel, message any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFormattedMessage", reflect.TypeOf((*MockISlack)(nil).AddFormattedMessage), channel, message)
}

// AddReaction mocks base method.
func (m *MockISlack) AddReaction(name string, item slack.MessageRef) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddReaction", name, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddReaction indicates an expected call of AddReaction.
func (mr *MockISlackMockRecorder) AddReaction(name, item any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddReaction", reflect.TypeOf((*MockISlack)(nil).AddReaction), name, item)
}

// RemoveReaction mocks base method.
func (m *MockISlack) RemoveReaction(name string, item slack.MessageRef) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveReaction", name, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveReaction indicates an expected call of RemoveReaction.
func (mr *MockISlackMockRecorder) RemoveReaction(name, item any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveReaction", reflect.TypeOf((*MockISlack)(nil).RemoveReaction), name, item)
}

// UploadFileWithContent mocks base method.
func (m *MockISlack) UploadFileWithContent(fileType, fileName, title, content string, messageRef slack.MessageRef) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFileWithContent", fileType, fileName, title, content, messageRef)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadFileWithContent indicates an expected call of UploadFileWithContent.
func (mr *MockISlackMockRecorder) UploadFileWithContent(fileType, fileName, title, content, messageRef any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFileWithContent", reflect.TypeOf((*MockISlack)(nil).UploadFileWithContent), fileType, fileName, title, content, messageRef)
}
