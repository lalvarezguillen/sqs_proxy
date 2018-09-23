package main

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func NewProxy(conf *AppConfig) *Proxy {
	fmt.Println("Config:", conf.Pretty())
	return &Proxy{
		Client: CreateSQSSession(),
		WG:     &sync.WaitGroup{},
		Conf:   conf,
	}
}

type Proxier interface {
	Start(*AppConfig)
	Hooker
}

type Proxy struct {
	Client SQSClientor
	WG     *sync.WaitGroup
	Conf   *AppConfig
	Hooker
	Mover
}

// Start the operations. Based on the proxy configuration,
// sets up a coroutine per source queue to handle the actual proxying.
func (p *Proxy) Start() {
	p.WG.Add(len(p.Conf.ProxyOps))
	for _, op := range p.Conf.ProxyOps {
		go p.Hook(&op, p.WG)
	}
	p.WG.Wait()
}

func CreateSQSSession() *sqs.SQS {
	sess := session.Must(session.NewSession())
	sqsSess := sqs.New(sess)
	return sqsSess
}

// SQSClientor is defined with the methods implemented by sqs.SQS, in order to be
// able to create structs that mock sqs.SQS
type SQSClientor interface {
	ReceiveMessage(i *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	SendMessage(i *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
	DeleteMessage(i *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
}
