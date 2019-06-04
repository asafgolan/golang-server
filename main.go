package main

import (
	"log"
	"net/http"
)

//var configuration Config

func main() {

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
