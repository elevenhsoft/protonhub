package main

import "net/http"

func Routing(mux *http.ServeMux) {
	mux.HandleFunc("GET /", RootHandler)
	mux.HandleFunc("GET /create", CreateHandler)
	mux.HandleFunc("POST /create", CreateDoneHandler)
	mux.HandleFunc("POST /run", RunHandler)
	mux.HandleFunc("GET /edit/{gameId}", EditHandler)
	mux.HandleFunc("POST /edit", EditDoneHandler)

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./static/assets/"))))
}
