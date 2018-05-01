package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	w.Write(data)
}

func main() {
	InitController()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Ending")
		CleanController()
		os.Exit(0)
	}()

	log.Println("Starting gpio API")

	r := mux.NewRouter()

	r.PathPrefix("/").
		Queries("gpio", "{gpio}").
		Queries("state", "{state}").
		Methods("GET").
		Handler(Controller{Index})

	r.PathPrefix("/").
		Methods("GET").
		HandlerFunc(Index)

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         ":4430",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("IP: %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
