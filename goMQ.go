package main

import (
	stomp "github.com/go-stomp/stomp"
)

type QueueClient struct {
	Addr string
}

func NewClient(addr string) *QueueClient {
	if addr == "" {
		addr = "localhost:61616"
	}
	return &QueueClient{addr}
}

func (q *QueueClient) Connect() (*stomp.Conn, error) {
	a, err := stomp.Dial("tcp", q.Addr)
	return a, err
}

func (q *QueueClient) Check() error {
	con, err := q.Connect()
	if err != nil {
		return err
	}
	con.Disconnect()
	return nil
}

func (q *QueueClient) Publish(targetQueue string, message []byte) error {
	conn, err := q.Connect()
	if err != nil {
		return err
	}
	defer conn.Disconnect()
	return conn.Send(targetQueue, "text/plain", message)
}

func (q *QueueClient) Subscribe(topic string, handler func(err error, message []byte)) error {
	con, err := q.Connect()
	if err != nil {
		return err
	}
	defer con.Disconnect()

	sub, err := con.Subscribe(topic, stomp.AckAuto)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		result := <-sub.C
		handler(result.Err, result.Body)
	}
}
