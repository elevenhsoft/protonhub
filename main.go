package main

import (
	"fmt"
	"log"
	"net/http"
)

const IP = "127.0.0.1"
const PORT = "8080"

func main() {

	mux := http.NewServeMux()
	Routing(mux)

	fmt.Printf("Starting server on http://%s:%s/\n", IP, PORT)

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", IP, PORT), mux)

	if err != nil {
		log.Fatalf("Server failed to start on port: %s", PORT)
	}
}
