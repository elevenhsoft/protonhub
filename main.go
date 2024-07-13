package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	Routing(mux)
	log.Fatal(http.ListenAndServe(":8080", mux))

}
