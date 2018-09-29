package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Mover outlines the interface that MessageMover should
// implement. Dummy implementations of Mover should be used
// in tests
type Mover interface {
	Move(*sqs.ReceiveMessageInput, TargetQueues) error
}

// MessagesMover moves messages between SQS queues.
type MessagesMover struct {
	Client SQSClientor
}

// Move reads some of the messages available in a source queue, and
// copies them to the destination queues, deleting them from the source queue
// afterwards.
func (p *MessagesMover) Move(src *sqs.ReceiveMessageInput, t TargetQueues) error {
	readResp, err := p.Client.ReceiveMessage(src)
	if err != nil {
		return err
	}
	if count := len(readResp.Messages); count > 0 {
		log.Println(fmt.Sprintf("%d messages to proxy from Queue %s",
			count, *src.QueueUrl))
	}

	// TODO: Look into batch writing and batch deleting
	for _, msg := range readResp.Messages {
		for _, q := range t {
			if err := p.moveMessage(msg, q); err != nil {
				return err
			}
		}
		if err := p.deleteMessage(msg, *src.QueueUrl); err != nil {
			return err
		}
	}
	return nil
}

// moveMessage inserts an sqs.Message to a new SQS Queue.
func (p *MessagesMover) moveMessage(m *sqs.Message, queueURL string) error {
	fmt.Println("called moveMovessage")
	writeParams := sqs.SendMessageInput{
		MessageBody: m.Body,
		QueueUrl:    aws.String(queueURL),
	}
	if _, err := p.Client.SendMessage(&writeParams); err != nil {
		return err
	}
	return nil
}

// deleteMessage removes an sqs.Message from its original queue
func (p *MessagesMover) deleteMessage(m *sqs.Message, queueURL string) error {
	fmt.Println("called deleteMessage")
	deleteParams := sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: m.ReceiptHandle,
	}
	if _, err := p.Client.DeleteMessage(&deleteParams); err != nil {
		return err
	}
	return nil
}
