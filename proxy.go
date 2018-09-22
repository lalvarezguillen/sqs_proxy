package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/urfave/cli"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "The path to the json file that contains the configuration to run this proxy",
		},
	}
	app.Action = func(c *cli.Context) error {
		conf := c.String("config")
		if conf == "" {
			return fmt.Errorf("--config is a required argument")
		}
		fmt.Println("Configuration file:", conf)

		proxyConf, err := loadConfig(conf)
		if err != nil {
			return err
		}
		StartProxy(proxyConf)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// StartProxy starts the operations. Based on the proxy configuration,
// sets up a coroutine per source queue to handle the actual proxying.
func StartProxy(conf *AppConfig) {
	// fmt.Println(fmt.Sprintf("Config:"))
	fmt.Println("Config:", conf.Pretty())
	s := createSQSSession()

	var wg sync.WaitGroup
	wg.Add(len(conf.ProxyOps))
	for _, op := range conf.ProxyOps {
		go HookToQueue(s, op, ProxyMessages, &wg)
	}
	wg.Wait()
}

func createSQSSession() *sqs.SQS {
	sess := session.Must(session.NewSession())
	sqsSess := sqs.New(sess)
	return sqsSess
}

// SQSClient is defined with the methods implemented by sqs.SQS, in order to be
// able to create structs that mock sqs.SQS
type SQSClient interface {
	ReceiveMessage(i *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	SendMessage(i *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
	DeleteMessage(i *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
}

// HookToQueue starts listening from a source queue, and handling the messages
// that come through.
func HookToQueue(s SQSClient, conf ProxySettings, p ProxyFunctionType, wg *sync.WaitGroup) error {
	defer wg.Done()
	readParams := sqs.ReceiveMessageInput{
		MaxNumberOfMessages: aws.Int64(10),
		QueueUrl:            aws.String(conf.Src),
		WaitTimeSeconds:     aws.Int64(20),
	}
	for {
		if err := p(s, &readParams, conf.Dest); err != nil {
			errIntro := fmt.Sprintf("Proxying from Queue %s has failed with error:", conf.Src)
			log.Println(errIntro, err)
			return err
		}
		time.Sleep(conf.Interval * time.Second)
	}
}

// ProxyFunctionType describes the kind of function that would take messages from
// a SQS queue and put them in a set of target queues.
type ProxyFunctionType func(SQSClient, *sqs.ReceiveMessageInput, []string) error

// ProxyMessages reads some of the messages available in a source queue, and
// copies them to the destination queues, deleting them from the source queue
// afterwards.
func ProxyMessages(s SQSClient, src *sqs.ReceiveMessageInput, dest []string) error {
	readResp, err := s.ReceiveMessage(src)
	if err != nil {
		return err
	}
	if count := len(readResp.Messages); count > 0 {
		log.Println(fmt.Sprintf("%d messages to proxy from Queue %s",
			count, *src.QueueUrl))
	}

	// TODO: Look into batch writing and batch deleting
	for _, msg := range readResp.Messages {
		for _, q := range dest {
			writeParams := sqs.SendMessageInput{
				MessageBody: msg.Body,
				QueueUrl:    aws.String(q),
			}
			if _, err := s.SendMessage(&writeParams); err != nil {
				return err
			}
		}
		deleteParams := sqs.DeleteMessageInput{
			QueueUrl:      src.QueueUrl,
			ReceiptHandle: msg.ReceiptHandle,
		}
		if _, err := s.DeleteMessage(&deleteParams); err != nil {
			return err
		}
	}
	return nil
}
