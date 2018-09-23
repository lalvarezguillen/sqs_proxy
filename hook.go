package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Hooker interface {
	Hook(*ProxySettings, *sync.WaitGroup)
	Mover
}

type QueueHook struct {
	Client SQSClientor
	Mover
}

// Hook starts listening from a source queue, and handling the messages
// that come through.
func (q *QueueHook) Hook(conf *ProxySettings, wg *sync.WaitGroup) error {
	defer wg.Done()
	readParams := sqs.ReceiveMessageInput{
		MaxNumberOfMessages: aws.Int64(10),
		QueueUrl:            aws.String(conf.Src),
		WaitTimeSeconds:     aws.Int64(20),
	}
	for {
		if err := q.Move(&readParams, conf.Dest); err != nil {
			errIntro := fmt.Sprintf("Proxying from Queue %s has failed with error:", conf.Src)
			log.Println(errIntro, err)
			return err
		}
		time.Sleep(conf.Interval * time.Second)
	}
}
