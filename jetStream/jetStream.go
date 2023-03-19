package jetStream

import "github.com/nats-io/nats.go"

func JetStreamInit() (nats.JetStreamContext, error) {
	// Connect to NATS to default URL (nats://127.0.0.1:4222)
	nc, err := nats.Connect(nats.DefaultURL)
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
