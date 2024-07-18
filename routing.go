package main

import (
	"io/fs"
	"log"
	"net/http"
)

func Routing(mux *http.ServeMux) {
	mux.HandleFunc("GET /", RootHandler)
	mux.HandleFunc("GET /create", CreateHandler)
	mux.HandleFunc("POST /create", CreateDoneHandler)
	mux.HandleFunc("POST /delete", DeleteHandler)
	mux.HandleFunc("GET /run/{gameId}", RunHandler)
	mux.HandleFunc("GET /stop/{gameId}", StopProcessHandler)
	mux.HandleFunc("GET /running-games", RunningGamesHandler)
	mux.HandleFunc("GET /winetricks/{gameId}/{verbs}", RunWinetricksHandler)
	mux.HandleFunc("GET /edit/{gameId}", EditHandler)
	mux.HandleFunc("POST /edit", EditDoneHandler)

	assetsRoot, err := fs.Sub(files, "static/assets")

	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServerFS(assetsRoot)))
}
