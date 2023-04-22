package main

import (
	"fmt"
	"log"
	"smartread/server"
)

// Main function creates and runs server
func main() {
	fmt.Println("statr")
	server, err := server.New()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := server.ListenAndServer(); err != nil {
		log.Fatal(err.Error())
	}
}
