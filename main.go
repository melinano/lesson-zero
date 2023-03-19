package main

import (
	"fmt"
	"log"
)

func main() {
	log.Println("Starting...")

	js, err := JetStreamInit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("JetStream context: ", js)

	err = CreateStream(js)
	if err != nil {
		log.Fatal(err)
	}

	subscribeOrderings(js)
	publishOrderings(js)

	log.Println("Exit...")
}
