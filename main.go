package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const IP = "127.0.0.1"
const PORT = "8090"

var (
	//go:embed static/*
	files     embed.FS
	templates map[string]*template.Template
)

func main() {
	// initalize store directory
	initStore()
	// initalize database
	conn := DbConnection()
	LaunchersTableInit(conn)

	// serve mux
	mux := http.NewServeMux()
	// set routings
	Routing(mux)

	// load templates to cache
	err := LoadTemplates()
	if err != nil {
		log.Fatal(err)
	}

	// start server
	fmt.Printf("Starting server on http://%s:%s/\n", IP, PORT)
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", IP, PORT), mux)
	if err != nil {
		log.Fatalf("Server failed to start on port: %s", PORT)
	}
}
