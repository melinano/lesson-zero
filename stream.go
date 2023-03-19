package main

import (
	"github.com/nats-io/nats.go"
	"log"
)

const (
	StreamName     = "MESSAGES"
	StreamSubjects = "MESSAGES.*"
)

func CreateStream(jetStream nats.JetStreamContext) error {
	stream, err := jetStream.StreamInfo(StreamName)

	// if stream not found, create it
	if stream == nil {
		log.Printf("Creating stream: %s\n", StreamName)

		_, err = jetStream.AddStream(&nats.StreamConfig{
			Name:     StreamName,
			Subjects: []string{StreamSubjects},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
