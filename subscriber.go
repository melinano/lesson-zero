package main

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/melinano/lesson-zero/models"
	"github.com/nats-io/nats.go"
	"log"
)

func subscribeOrderings(js nats.JetStreamContext, pool *pgxpool.Pool, orderingsMap *map[string]models.Ordering) {
	_, err := js.Subscribe("MESSAGES.*", func(msg *nats.Msg) {
		err := msg.Ack()
		if err != nil {
			log.Println("Unable to Ack", err)
			return
		}
		var ordering models.Ordering
		err = json.Unmarshal(msg.Data, &ordering)
		if err != nil {
			log.Println(err)
		}
		// save the retrieved Ordering into the map
		(*orderingsMap)[ordering.OrderUid] = ordering
		// insert Ordering into Postgres DB
		err = insertOrderingIntoDB(pool, ordering)
		if err != nil {
			log.Printf("%s\nOrdering(order_uid: %s) skipped!", err, ordering.OrderUid)
		}

		log.Printf("Subscriber => Subject: %s - ID: %s", msg.Subject, ordering.OrderUid)
	}, nats.Durable("MESSAGE"))

	if err != nil {
		fmt.Println("Error subscribing to JetStream subject:", err)
	}

	select {}
}
