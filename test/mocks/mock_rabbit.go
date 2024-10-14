package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRabbitMQClient is a mock of RabbitMQClient interface.
type MockRabbitMQClient struct {
	ctrl     *gomock.Controller
	recorder *MockRabbitMQClientMockRecorder
}

// MockRabbitMQClientMockRecorder is the mock recorder for MockRabbitMQClient.
type MockRabbitMQClientMockRecorder struct {
	mock *MockRabbitMQClient
}

// NewMockRabbitMQClient creates a new mock instance.
func NewMockRabbitMQClient(ctrl *gomock.Controller) *MockRabbitMQClient {
	mock := &MockRabbitMQClient{ctrl: ctrl}
	mock.recorder = &MockRabbitMQClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRabbitMQClient) EXPECT() *MockRabbitMQClientMockRecorder {
	return m.recorder
}

// ConsumeMessage mocks base method.
func (m *MockRabbitMQClient) ConsumeMessage(queue string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConsumeMessage", queue)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConsumeMessage indicates an expected call of ConsumeMessage.
func (mr *MockRabbitMQClientMockRecorder) ConsumeMessage(queue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConsumeMessage", reflect.TypeOf((*MockRabbitMQClient)(nil).ConsumeMessage), queue)
}
