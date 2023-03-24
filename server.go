package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// establish a http server on localhost and port 8080
func startHTTPServer() {
	http.HandleFunc("/orderings", getOrderingByIDHandler)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func getOrderingByIDHandler(w http.ResponseWriter, r *http.Request) {
	// query the request url
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	var orderingJSON []byte // Json byte slice
	var err error
	// check if key exists in map
	if ordering, ok := orderingsMap[id]; ok {
		orderingJSON, err = json.Marshal(ordering)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// write to server
	w.Write(orderingJSON)
	return
}
