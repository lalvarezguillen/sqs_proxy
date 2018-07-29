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
		fmt.Println("Using config file: " + conf)

		proxyConf, err := loadConfig(conf)
		if err != nil {
			return err
		}
		startProxy(proxyConf)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func startProxy(conf *AppConfig) {
	fmt.Println(fmt.Sprintf("starting proxy with conf %+v", conf))
	s := createSQSSession()

	var wg sync.WaitGroup
	wg.Add(len(conf.ProxyOps))
	for _, op := range conf.ProxyOps {
		go hookToQueue(s, &op, &wg)
	}
	wg.Wait()
}

func createSQSSession() *sqs.SQS {
	sess := session.Must(session.NewSession())
	sqsSess := sqs.New(sess)
	return sqsSess
}

type SQSClient interface {
	ReceiveMessage(i *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	SendMessage(i *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
	DeleteMessage(i *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
}

func hookToQueue(s SQSClient, conf *ProxySettings, wg *sync.WaitGroup) {
	defer wg.Done()
	readParams := sqs.ReceiveMessageInput{
		MaxNumberOfMessages: aws.Int64(10),
		QueueUrl:            aws.String(conf.Src),
		WaitTimeSeconds:     aws.Int64(20),
	}
	for {
		proxyMessages(s, &readParams, conf.Dest)
		time.Sleep(conf.Interval * time.Second)
	}
}

func proxyMessages(s SQSClient, src *sqs.ReceiveMessageInput, dest []string) {
	readResp, err := s.ReceiveMessage(src)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("%d messages to proxy from Queue %s",
		len(readResp.Messages), *src.QueueUrl))
	// TODO: Look into batch writing and batch deleting
	for _, msg := range readResp.Messages {
		for _, q := range dest {
			writeParams := sqs.SendMessageInput{
				MessageBody: msg.Body,
				QueueUrl:    aws.String(q),
			}
			if _, err := s.SendMessage(&writeParams); err != nil {
				panic(err)
			}
		}
		deleteParams := sqs.DeleteMessageInput{
			QueueUrl:      src.QueueUrl,
			ReceiptHandle: msg.ReceiptHandle,
		}
		if _, err := s.DeleteMessage(&deleteParams); err != nil {
			panic(err)
		}
	}
}
