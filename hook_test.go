package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
)

func TestHookToQueueError(t *testing.T) {
	// Setup
	var proxyFuncInvocations int
	var proxyFuncSQSClient SQSClientor
	var proxyFuncReceiveMessageInput *sqs.ReceiveMessageInput
	var proxyFuncDestQueues []string
	dummyProxyFunc := func(s SQSClientor, src *sqs.ReceiveMessageInput, dest []string) error {
		proxyFuncInvocations++
		proxyFuncSQSClient = s
		proxyFuncReceiveMessageInput = src
		proxyFuncDestQueues = dest
		return fmt.Errorf("Dummy error")
	}

	c := &MockedSQS{}
	conf := ProxySettings{
		Src:  "source-queue",
		Dest: []string{"target-queue-1", "target-queue-2"},
	}
	var wg sync.WaitGroup
	wg.Add(1)

	// Actual tests
	assert.Error(t, HookToQueue(c, conf, dummyProxyFunc, &wg))
	assert.Equal(t, c, proxyFuncSQSClient)
	assert.Equal(t, conf.Dest, proxyFuncDestQueues)
	assert.Equal(t, conf.Src, *proxyFuncReceiveMessageInput.QueueUrl)
}
