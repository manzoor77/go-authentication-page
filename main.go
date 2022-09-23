package main

import (
	"log"
	"net/http"

	"loginPackage.com/handlers"
)

func main() {
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/home", handlers.Home)
	//http.HandleFunc("/refresh", refresh)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
