package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	proxyOps, err := loadProxySettings(config.ProxySettingsFile)
	if err != nil {
		panic(err)
	}

	s, err := createSQSSession(config)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(proxyOps))
	for _, op := range proxyOps {
		go proxyQueue(s, op.Src, op.Dest, op.Interval, &wg)
	}
	wg.Wait()
}

func createSQSSession(c *AppConfig) (*sqs.SQS, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(c.AWSRegion),
	}))
	sqsSess := sqs.New(sess)
	return sqsSess, nil
}

func proxyQueue(s *sqs.SQS, srcQ string, destQs []string, interval time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	readParams := sqs.ReceiveMessageInput{
		MaxNumberOfMessages: aws.Int64(10),
		QueueUrl:            aws.String(srcQ),
		WaitTimeSeconds:     aws.Int64(20),
	}
	for {
		readResp, err := s.ReceiveMessage(&readParams)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("%d messages to proxy from Queue %s", len(readResp.Messages), srcQ))
		// TODO: Look into batch writing and batch deleting
		for _, msg := range readResp.Messages {
			for _, q := range destQs {
				writeParams := sqs.SendMessageInput{
					MessageBody: msg.Body,
					QueueUrl:    aws.String(q),
				}
				if _, err := s.SendMessage(&writeParams); err != nil {
					panic(err)
				}
			}
			deleteParams := sqs.DeleteMessageInput{
				QueueUrl:      aws.String(srcQ),
				ReceiptHandle: msg.ReceiptHandle,
			}
			if _, err := s.DeleteMessage(&deleteParams); err != nil {
				panic(err)
			}
		}
		time.Sleep(interval * time.Second)
	}
}
