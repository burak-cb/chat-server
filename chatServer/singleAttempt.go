package main

/*
import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	numberOfClients := 0
	sliceOfClients := make(map[net.Conn]int)

	newConnections := make(chan net.Conn)
	disconnectedConnections := make(chan net.Conn)

	incomingMessages := make(chan string)


	tcpServer, serverError := net.Listen("tcp", ":8080")
	if serverError != nil {
		fmt.Println("Connection Error:", serverError)
		os.Exit(1)
	}

	go func() {
		for {
			incomingConnection, connectionError := tcpServer.Accept()
			if connectionError != nil {
				fmt.Println("Incoming Connection Error:", connectionError)
				os.Exit(1)
			}

			newConnections <- incomingConnection
		}
	}()

	for {
		select {
		case newClient := <-newConnections:
			log.Println("New Client:", numberOfClients)

			sliceOfClients[newClient] = numberOfClients
			numberOfClients += 1

			go func(newConnection net.Conn, clientID int) {
				inputReader := bufio.NewReader(newConnection)

				for {
					incomingMessage, messageError := inputReader.ReadString('\n')
					if messageError != nil {
						break
					}

					incomingMessages <- fmt.Sprintf("Client #%d: %s", clientID, incomingMessage)
				}

				// If error occurs during the reading.
				disconnectedConnections <- newConnection
			}(newClient, sliceOfClients[newClient])

		case incomingMessage := <-incomingMessages:
			for eachClient, _ := range sliceOfClients {
				go func(clientConnection net.Conn, incomingMessage string) {
					_, messageError := clientConnection.Write([]byte(incomingMessage))

					if messageError != nil {
						disconnectedConnections <- clientConnection
					}
				}(eachClient, incomingMessage)
			}
			log.Println(incomingMessage)

		case disconnectedConnection := <-disconnectedConnections:
			log.Println("Disconnection: Client ", sliceOfClients[disconnectedConnection])
			delete(sliceOfClients, disconnectedConnection)
		}
	}

}
*/
