package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
)

func TestProxyMessages(t *testing.T) {
	// Setup
	c := &MockedSQS{}
	target := TargetQueues{
		"http://queues.com/dummy-destination",
	}

	src := sqs.ReceiveMessageInput{QueueUrl: aws.String("http://queues.com/dummy")}
	outp := sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{
			&sqs.Message{Body: aws.String("dummy message 1"), ReceiptHandle: aws.String("dummy-1")},
			&sqs.Message{Body: aws.String("dummy message 2"), ReceiptHandle: aws.String("dummy-2")},
		},
	}
	c.On("ReceiveMessage", &src).Return(outp, nil)

	smInput1 := sqs.SendMessageInput{
		MessageBody: aws.String("dummy message 1"),
		QueueUrl:    aws.String(target[0]),
	}
	c.On("SendMessage", &smInput1).Return(&sqs.SendMessageOutput{}, nil)
	smInput2 := sqs.SendMessageInput{
		MessageBody: aws.String("dummy message 2"),
		QueueUrl:    aws.String(target[0]),
	}
	c.On("SendMessage", &smInput2).Return(&sqs.SendMessageOutput{}, nil)

	dmInput1 := sqs.DeleteMessageInput{
		QueueUrl:      src.QueueUrl,
		ReceiptHandle: aws.String("dummy-1"),
	}
	c.On("DeleteMessage", &dmInput1).Return(&sqs.DeleteMessageOutput{}, nil)
	dmInput2 := sqs.DeleteMessageInput{
		QueueUrl:      src.QueueUrl,
		ReceiptHandle: aws.String("dummy-2"),
	}
	c.On("DeleteMessage", &dmInput2).Return(&sqs.DeleteMessageOutput{}, nil)

	// Actual test
	m := MessagesMover{Client: c}
	assert.NoError(t, m.Move(&src, target))
	c.AssertExpectations(t)
	c.AssertNumberOfCalls(t, "ReceiveMessage", 1)
	c.AssertNumberOfCalls(t, "SendMessage", 2)
	c.AssertNumberOfCalls(t, "DeleteMessage", 2)
}

func TestNoMessagesToProxy(t *testing.T) {
	// Setup
	c := &MockedSQS{}
	target := TargetQueues{
		"http://queues.com/dummy-destination",
	}

	src := sqs.ReceiveMessageInput{QueueUrl: aws.String("http://queues.com/dummy")}
	outp := sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{},
	}
	c.On("ReceiveMessage", &src).Return(outp, nil)

	smInput1 := sqs.SendMessageInput{
		MessageBody: aws.String("dummy message 1"),
		QueueUrl:    aws.String(target[0]),
	}
	c.On("SendMessage", &smInput1).Return(&sqs.SendMessageOutput{}, nil)
	smInput2 := sqs.SendMessageInput{
		MessageBody: aws.String("dummy message 2"),
		QueueUrl:    aws.String(target[0]),
	}
	c.On("SendMessage", &smInput2).Return(&sqs.SendMessageOutput{}, nil)

	dmInput1 := sqs.DeleteMessageInput{
		QueueUrl:      src.QueueUrl,
		ReceiptHandle: aws.String("dummy-1"),
	}
	c.On("DeleteMessage", &dmInput1).Return(&sqs.DeleteMessageOutput{}, nil)
	dmInput2 := sqs.DeleteMessageInput{
		QueueUrl:      src.QueueUrl,
		ReceiptHandle: aws.String("dummy-2"),
	}
	c.On("DeleteMessage", &dmInput2).Return(&sqs.DeleteMessageOutput{}, nil)

	// Actual test
	m := MessagesMover{Client: c}
	assert.NoError(t, m.Move(&src, target))
	c.AssertNumberOfCalls(t, "ReceiveMessage", 1)
	c.AssertNotCalled(t, "SendMessage")
	c.AssertNotCalled(t, "DeleteMessage")
}

func TestProxyMessagesErrorReading(t *testing.T) {
	// Setup
	c := &MockedSQS{}
	target := TargetQueues{
		"http://queues.com/dummy-destination",
	}

	src := sqs.ReceiveMessageInput{QueueUrl: aws.String("http://queues.com/dummy")}
	outp := sqs.ReceiveMessageOutput{Messages: []*sqs.Message{}}
	c.On("ReceiveMessage", &src).Return(outp, fmt.Errorf("Reading Failed"))
	// Actual test
	m := MessagesMover{Client: c}
	assert.Error(t, m.Move(&src, target))
	c.AssertExpectations(t)
	c.AssertNumberOfCalls(t, "ReceiveMessage", 1)
	c.AssertNotCalled(t, "SendMessage")
	c.AssertNotCalled(t, "DeleteMessage")
}

func TestProxyMessagesErrorSending(t *testing.T) {
	// Setup
	c := &MockedSQS{}
	target := TargetQueues{
		"http://queues.com/dummy-destination",
	}

	src := sqs.ReceiveMessageInput{QueueUrl: aws.String("http://queues.com/dummy")}
	outp := sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{
			&sqs.Message{Body: aws.String("dummy message 1"), ReceiptHandle: aws.String("dummy-1")},
			&sqs.Message{Body: aws.String("dummy message 2"), ReceiptHandle: aws.String("dummy-2")},
		},
	}
	c.On("ReceiveMessage", &src).Return(outp, nil)

	smInput1 := sqs.SendMessageInput{
		MessageBody: aws.String("dummy message 1"),
		QueueUrl:    aws.String(target[0]),
	}
	c.On("SendMessage", &smInput1).Return(&sqs.SendMessageOutput{}, fmt.Errorf("Dummy Error"))

	// Actual test
	m := MessagesMover{Client: c}
	assert.Error(t, m.Move(&src, target))
	c.AssertExpectations(t)
	c.AssertNumberOfCalls(t, "ReceiveMessage", 1)
	c.AssertNumberOfCalls(t, "SendMessage", 1)
	c.AssertNotCalled(t, "DeleteMessage")
}

func TestProxyMessagesErrorDeleting(t *testing.T) {
	// Setup
	c := &MockedSQS{}
	target := TargetQueues{
		"http://queues.com/dummy-destination",
	}

	src := sqs.ReceiveMessageInput{QueueUrl: aws.String("http://queues.com/dummy")}
	outp := sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{
			&sqs.Message{Body: aws.String("dummy message 1"), ReceiptHandle: aws.String("dummy-1")},
			&sqs.Message{Body: aws.String("dummy message 2"), ReceiptHandle: aws.String("dummy-2")},
		},
	}
	c.On("ReceiveMessage", &src).Return(outp, nil)

	smInput1 := sqs.SendMessageInput{
		MessageBody: aws.String("dummy message 1"),
		QueueUrl:    aws.String(target[0]),
	}
	c.On("SendMessage", &smInput1).Return(&sqs.SendMessageOutput{}, nil)

	dmInput1 := sqs.DeleteMessageInput{
		QueueUrl:      src.QueueUrl,
		ReceiptHandle: aws.String("dummy-1"),
	}
	c.On("DeleteMessage", &dmInput1).Return(&sqs.DeleteMessageOutput{}, fmt.Errorf("Dummy Error"))

	// Actual test
	m := MessagesMover{Client: c}
	assert.Error(t, m.Move(&src, target))
	c.AssertExpectations(t)
	c.AssertNumberOfCalls(t, "ReceiveMessage", 1)
	c.AssertNumberOfCalls(t, "SendMessage", 1)
	c.AssertNumberOfCalls(t, "DeleteMessage", 1)
}
