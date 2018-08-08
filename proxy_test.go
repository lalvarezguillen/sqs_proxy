package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
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

func TestProxyMessages(t *testing.T) {
	// Setup
	c := &MockedSQS{}
	d := []string{
		"http://queues.com/dummy-destination",
	}

	i := sqs.ReceiveMessageInput{QueueUrl: aws.String("http://queues.com/dummy")}
	outp := sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{
			&sqs.Message{Body: aws.String("dummy message 1"), ReceiptHandle: aws.String("dummy-1")},
			&sqs.Message{Body: aws.String("dummy message 2"), ReceiptHandle: aws.String("dummy-2")},
		},
	}
	c.On("ReceiveMessage", &i).Return(outp, nil)

	smInput1 := sqs.SendMessageInput{
		MessageBody: aws.String("dummy message 1"),
		QueueUrl:    aws.String(d[0]),
	}
	c.On("SendMessage", &smInput1).Return(&sqs.SendMessageOutput{}, nil)
	smInput2 := sqs.SendMessageInput{
		MessageBody: aws.String("dummy message 2"),
		QueueUrl:    aws.String(d[0]),
	}
	c.On("SendMessage", &smInput2).Return(&sqs.SendMessageOutput{}, nil)

	dmInput1 := sqs.DeleteMessageInput{
		QueueUrl:      i.QueueUrl,
		ReceiptHandle: aws.String("dummy-1"),
	}
	c.On("DeleteMessage", &dmInput1).Return(&sqs.DeleteMessageOutput{}, nil)
	dmInput2 := sqs.DeleteMessageInput{
		QueueUrl:      i.QueueUrl,
		ReceiptHandle: aws.String("dummy-2"),
	}
	c.On("DeleteMessage", &dmInput2).Return(&sqs.DeleteMessageOutput{}, nil)

	// Actual test
	ProxyMessages(c, &i, d)
	c.AssertExpectations(t)
	c.AssertNumberOfCalls(t, "ReceiveMessage", 1)
	c.AssertNumberOfCalls(t, "SendMessage", 2)
	c.AssertNumberOfCalls(t, "DeleteMessage", 2)
}

func TestNoMessagesToProxy(t *testing.T) {
	// Setup
	c := &MockedSQS{}
	d := []string{
		"http://queues.com/dummy-destination",
	}

	i := sqs.ReceiveMessageInput{QueueUrl: aws.String("http://queues.com/dummy")}
	outp := sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{},
	}
	c.On("ReceiveMessage", &i).Return(outp, nil)

	smInput1 := sqs.SendMessageInput{
		MessageBody: aws.String("dummy message 1"),
		QueueUrl:    aws.String(d[0]),
	}
	c.On("SendMessage", &smInput1).Return(&sqs.SendMessageOutput{}, nil)
	smInput2 := sqs.SendMessageInput{
		MessageBody: aws.String("dummy message 2"),
		QueueUrl:    aws.String(d[0]),
	}
	c.On("SendMessage", &smInput2).Return(&sqs.SendMessageOutput{}, nil)

	dmInput1 := sqs.DeleteMessageInput{
		QueueUrl:      i.QueueUrl,
		ReceiptHandle: aws.String("dummy-1"),
	}
	c.On("DeleteMessage", &dmInput1).Return(&sqs.DeleteMessageOutput{}, nil)
	dmInput2 := sqs.DeleteMessageInput{
		QueueUrl:      i.QueueUrl,
		ReceiptHandle: aws.String("dummy-2"),
	}
	c.On("DeleteMessage", &dmInput2).Return(&sqs.DeleteMessageOutput{}, nil)

	// Actual test
	ProxyMessages(c, &i, d)
	c.AssertNumberOfCalls(t, "ReceiveMessage", 1)
	c.AssertNotCalled(t, "SendMessage")
	c.AssertNotCalled(t, "DeleteMessage")
}
