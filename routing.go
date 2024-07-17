package main

import "net/http"

func Routing(mux *http.ServeMux) {
	mux.HandleFunc("GET /", RootHandler)
	mux.HandleFunc("GET /create", CreateHandler)
	mux.HandleFunc("POST /create", CreateDoneHandler)
	mux.HandleFunc("GET /run/{gameId}", RunHandler)
	mux.HandleFunc("POST /create-lock", CreateProcessLockHandler)
	mux.HandleFunc("POST /remove-lock", RemoveProcessLockHandler)
	mux.HandleFunc("GET /winetricks/{gameId}/{verbs}", RunWinetricksHandler)
	mux.HandleFunc("GET /edit/{gameId}", EditHandler)
	mux.HandleFunc("POST /edit", EditDoneHandler)

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./static/assets/"))))
}
