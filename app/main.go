package main

import (
 	"log"
 	"net/http"
 	"github.com/gorilla/mux"
)

var	myRoutes = map[string]func(http.ResponseWriter, *http.Request){
		"/": Index,
		"/meta": MyMeta,
		"/servers": MyPeers}

func main() {

	// Create the HTTP router
	router := mux.NewRouter().StrictSlash(true);
	for key, value := range myRoutes {
		router.HandleFunc(key, value);
	}	

	// Start serving requests
	log.Fatal(http.ListenAndServe(":8080", router))
}