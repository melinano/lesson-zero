package main

import (
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/melinano/lesson-zero/models"
	"log"
)

// Go-Map as in-memory storage
var orderingsMap = make(map[string]models.Ordering)

func main() {
	log.Println("Starting...")
	// opening DB as pool of connections
	var pool *pgxpool.Pool
	pool = startDB()
	defer pool.Close()

	log.Println("Fetching data from postgres...")
	// fetch data from DB into memory
	fetchOrderings(pool, &orderingsMap)

	// starting HTTP Server in a goroutine
	go startHTTPServer()

	// initiate JetStream
	js, err := JetStreamInit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("JetStream context: ", js)

	// Add a stream
	err = CreateStream(js)
	if err != nil {
		log.Fatal(err)
	}

	// start a subscription in a goroutine
	go subscribeOrderings(js, pool, &orderingsMap)
	// start a publisher in a goroutine
	go publishOrderings(js)

	select {}
}
