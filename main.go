package main

import (
	"log"
	"net/http"

	"loginPackage.com/handlers"
)

func main() {
	//create route handle function here
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/home", handlers.Home)
	http.HandleFunc("/refresh", handlers.Refresh)
	// set port 8080 for server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
