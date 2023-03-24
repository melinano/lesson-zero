package main

import (
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/melinano/lesson-zero/models"
	"log"
)

var orderingsMap = make(map[string]models.Ordering)

func main() {
	log.Println("Starting...")
	// opening DB
	var pool *pgxpool.Pool
	pool = startDB()
	defer pool.Close()
	log.Println("Fetching data from postgres...")
	// create map as in-memory storage

	// fetch data from DB into memory
	fetchOrderings(pool, &orderingsMap)

	// starting HTTP Server
	go startHTTPServer()

	js, err := JetStreamInit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("JetStream context: ", js)

	err = CreateStream(js)
	if err != nil {
		log.Fatal(err)
	}

	go subscribeOrderings(js, pool, &orderingsMap)
	go publishOrderings(js)

	select {}
	log.Println("Exit...")
}
