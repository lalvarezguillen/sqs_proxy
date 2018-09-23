package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MoverMock struct {
	mock.Mock
}

func (m *MoverMock) Move(i *sqs.ReceiveMessageInput, t TargetQueues) error {
	args := m.Called(i, t)
	return args.Error(0)
}

func TestHookToQueueError(t *testing.T) {
	// Setup
	targ := TargetQueues{"https://queues.com/dummy-dest"}
	conf := ProxySettings{
		Src:  "https://queues.com/dummy-src",
		Dest: targ,
	}
	src := sqs.ReceiveMessageInput{
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(20),
		QueueUrl:            aws.String(conf.Src),
	}

	m := MoverMock{}
	m.On("Move", &src, targ).Return(fmt.Errorf("Dummy Error"))

	h := QueueHook{
		Mover:  &m,
		Client: &MockedSQS{},
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// Actual tests
	assert.Error(t, h.Hook(&conf, &wg))
	m.AssertExpectations(t)
	m.AssertCalled(t, "Move", &src, targ)
}
