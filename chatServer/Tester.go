package main

import (
	"fmt"
	"log"
	"net"
)

const portNumber string = "localhost:8080"

func main() {
	connectionListener, connectionError := net.Listen("tcp", portNumber)

	// ---------------------------------------------
	if connectionError != nil {
		fmt.Println("Problem!")
	}
	// ---------------------------------------------

	hostAddress, portNumber, _ := net.SplitHostPort(connectionListener.Addr().String())

	defer func(connectionListener net.Listener) {
		closureError := connectionListener.Close()
		if closureError != nil {
			fmt.Println("Error CLosing the Connection.")
		}
	}(connectionListener)

	log.Printf("Listening on Host: %s - Port: %s\n", hostAddress, portNumber)

	// ----------------------------------------------------------------------------------------------------------------
	for {
		// Listening for incoming connections.
		incomingConnection, newConnectionError := connectionListener.Accept()

		// ---------------------------------------------

		if newConnectionError != nil {
			panic(newConnectionError)
		}
		// ---------------------------------------------

		// Handle each new connection on a separate go routine.
		go func(eachConnection net.Conn) {
			messageBuffer := make([]byte, 1024)
			messageLength, messageError := incomingConnection.Read(messageBuffer)

			// ---------------------------------------------
			if messageError != nil {
				fmt.Printf("Error Reading: %#v\n", messageError)
				return
			}
			// ---------------------------------------------

			fmt.Printf("Message Received: %s\n", string(messageBuffer[:messageLength]))

			incomingConnection.Write([]byte("Message Received.\n"))
			incomingConnection.Close()
		}(incomingConnection)
		// ----------------------------------------------------------------------------------------------------------------
	}
}
