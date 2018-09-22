package main

import (
	"github.com/aws/aws-sdk-go/service/sqs"
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
