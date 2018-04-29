package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func Control(w http.ResponseWriter, r *http.Request) {
	gpio := r.URL.Query().Get("gpio")
	fmt.Fprintf(w, "Your Are controlling the gpio %s", gpio)
}

func main() {
	log.Println("Starting gpio API")

	r := mux.NewRouter()

	r.PathPrefix("/").
		Queries("gpio", "{gpio}").
		Queries("state", "{state}").
		Methods("GET").
		HandlerFunc(Control)

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
