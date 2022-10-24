package main

import (
	"log"
	"net/http"
)

func main() {
	multiplexerRouter := routes()

	log.Println("Starting a Web Server on Port 8080.")
	_ = http.ListenAndServe(":8080", multiplexerRouter)
}
