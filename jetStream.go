package main

import (
	"github.com/nats-io/nats.go"
	"time"
)

func JetStreamInit() (nats.JetStreamContext, error) {
	// Connect to NATS to default URL (nats://127.0.0.1:4222)
	opts := nats.Options{
		Url:     nats.DefaultURL,
		Timeout: 5 * time.Second,
	}

	nc, err := opts.Connect()
	if err != nil {
		return nil, err
	}

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	return js, nil
}
