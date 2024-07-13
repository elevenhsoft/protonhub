package main

import "net/http"

func Routing(mux *http.ServeMux) {
	mux.HandleFunc("GET /", RootHandler)
	mux.HandleFunc("GET /create", CreateHandler)
	mux.HandleFunc("POST /create", CreateDoneHandler)
}
