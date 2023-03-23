package main

import (
	"encoding/json"
	"fmt"
	"github.com/melinano/lesson-zero/models"
	"github.com/nats-io/nats.go"
	"log"
)

func subscribeOrderings(js nats.JetStreamContext) {
	_, err := js.Subscribe("MESSAGES.*", func(msg *nats.Msg) {
		err := msg.Ack()
		if err != nil {
			log.Println("Unable to Ack", err)
			return
		}
		var ordering models.Ordering
		err = json.Unmarshal(msg.Data, &ordering)
		if err != nil {
			log.Fatal(err)
		}

		startDB(ordering)

		log.Printf("Subscriber => Subject: %s - ID: %s", msg.Subject, ordering.OrderUid)
	}, nats.Durable("MESSAGE"))

	if err != nil {
		fmt.Println("Error subscribing to JetStream subject:", err)
		return
	}
}
