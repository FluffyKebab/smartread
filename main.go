package main

import (
	"log"
	"smartread/server"
)

// Main function creates and runs server
func main() {
	server, err := server.New()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := server.ListenAndServer(); err != nil {
		log.Fatal(err.Error())
	}
}
