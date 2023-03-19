package pub

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

func getOrderings() ([]Ordering, error) {
	rawOrderings, _ := ioutil.ReadFile("./model.json")
	var orderingObj []Ordering
	err := json.Unmarshal(rawOrderings, &orderingObj)

	return orderingObj, err
}

func publishOrderings(js nats.JetStreamContext) {
	orderings, err := getOrderings()
	if err != nil {
		log.Println(err)
		return
	}

	for _, oneOrdering := range orderings {

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
			log.Printf("Publisher => Message:%s\n", oneOrdering.Uid)
		}
	}
}
