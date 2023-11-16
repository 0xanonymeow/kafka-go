// Code generated by MockGen. DO NOT EDIT.
// Source: client.go
//
// Generated by this command:
//
//	mockgen -source=client.go -package=mockclient -destination=testings/mock_client/mockclient.go
//
// Package mockclient is a generated GoMock package.
package mockclient

import (
	reflect "reflect"

	consumer "github.com/0xanonymeow/kafka-go/consumer"
	message "github.com/0xanonymeow/kafka-go/message"
	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Produce mocks base method.
func (m *MockClient) Produce(arg0 message.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Produce", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Produce indicates an expected call of Produce.
func (mr *MockClientMockRecorder) Produce(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Produce", reflect.TypeOf((*MockClient)(nil).Produce), arg0)
}

// RegisterConsumerHandler mocks base method.
func (m *MockClient) RegisterConsumerHandler(arg0 consumer.Topic) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterConsumerHandler", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterConsumerHandler indicates an expected call of RegisterConsumerHandler.
func (mr *MockClientMockRecorder) RegisterConsumerHandler(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterConsumerHandler", reflect.TypeOf((*MockClient)(nil).RegisterConsumerHandler), arg0)
}
