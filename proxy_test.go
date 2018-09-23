package main

import (
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedSQS struct {
	mock.Mock
}

func (c *MockedSQS) ReceiveMessage(i *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	args := c.Called(i)
	ro, ok := args.Get(0).(sqs.ReceiveMessageOutput)
	if !ok {
		panic("Failed to cast to *sqs.ReceiveMessageOutput")
	}
	return &ro, args.Error(1)
}

func (c *MockedSQS) SendMessage(i *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	args := c.Called(i)
	so, ok := args.Get(0).(*sqs.SendMessageOutput)
	if !ok {
		panic("Failed to cast to *sqs.SendMessageOutput")
	}
	return so, args.Error(1)
}

func (c *MockedSQS) DeleteMessage(i *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	args := c.Called(i)
	do, ok := args.Get(0).(*sqs.DeleteMessageOutput)
	if !ok {
		panic("Failed to cast to *sqs.DeleteMessageOutput")
	}
	return do, args.Error(1)
}

type MockedHook struct {
	HookCounter int
}

func (h *MockedHook) Hook(p *ProxySettings, wg *sync.WaitGroup) {
	defer wg.Done()
	h.HookCounter++
}

func (h *MockedHook) Move(i *sqs.ReceiveMessageInput, t TargetQueues) error {
	return nil
}

func TestProxyStart(t *testing.T) {
	conf := AppConfig{
		ProxyOps: []ProxySettings{
			ProxySettings{
				Src:      "https://queues.com/dummy-src",
				Dest:     TargetQueues{"https://queues.com/dummy-dest"},
				Interval: time.Duration(10),
			},
			ProxySettings{
				Src:      "https://queues.com/dummy-src-2",
				Dest:     TargetQueues{"https://queues.com/dummy-dest-2"},
				Interval: time.Duration(10),
			},
		},
	}
	h := MockedHook{}
	p := Proxy{
		WG:     &sync.WaitGroup{},
		Conf:   &conf,
		Hooker: &h,
	}
	p.Start()
	assert.Equal(t, 2, h.HookCounter)
}

func TestCreateSQSSession(t *testing.T) {
	assert.NotPanics(t, func() { CreateSQSSession() })
}

func TestNewProxy(t *testing.T) {
	conf := AppConfig{
		ProxyOps: []ProxySettings{
			ProxySettings{
				Src:      "https://queues.com/dummy-src",
				Dest:     TargetQueues{"https://queues.com/dummy-dest"},
				Interval: time.Duration(10),
			},
		},
	}

	proxy := NewProxy(&conf)
	assert.Equal(t, &conf, proxy.Conf)
}
