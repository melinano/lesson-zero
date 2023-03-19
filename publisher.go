package main

import (
	"encoding/json"
	"github.com/melinano/lesson-zero/models"
	"github.com/nats-io/nats.go"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

func getOrderings() ([]models.Ordering, error) {
	rawOrderings, _ := ioutil.ReadFile("./model.json")
	var orderingObj models.Ordering
	err := json.Unmarshal(rawOrderings, &orderingObj)

	return []models.Ordering{orderingObj}, err
}

func publishOrderings(js nats.JetStreamContext) {
	orderings, err := getOrderings()
	if err != nil {
		log.Println(err)
		return
	}
	oneOrdering := orderings[0]

	//for _, oneOrdering := range orderings {
	for i := 0; i <= 30; i++ {

		// create random message intervals to slow down
		r := rand.Intn(1500)
		time.Sleep(time.Duration(r) * time.Millisecond)

		orderingString, err := json.Marshal(oneOrdering)
		if err != nil {
			log.Println(err)
			continue
		}

		// publish to MESSAGES.ordering subject
		_, err = js.Publish("MESSAGES.orderings", orderingString)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Publisher => Message:%s\n", oneOrdering.OrderUid)
		}
	}
}
