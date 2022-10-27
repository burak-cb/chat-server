package main

import (
	"net"
	"testing"
	"time"
)

var chatRoom ChatRoom
var chatUser ChatUser
var testChannel chan string

func init() {
	go StartServer()
	time.Sleep(time.Second) // Putting into the sleep works.
	// testChannel <- "Server is UP, starting \"TestTCPServerRunning\" test." // It it a better practice.
}

func TestTCPServerRunning(t *testing.T) {
	// fmt.Println(<-testChannel)
	type serverStruct struct {
		serverProtocol string
		portAddress    string
	}

	chatServer := serverStruct{
		serverProtocol: "tcp",
		portAddress:    ":8080",
	}

	testConnection, connectionError := net.Dial(chatServer.serverProtocol, chatServer.portAddress)
	if connectionError != nil {
		t.Error("ERROR Connecting to the Server: ", connectionError)
	}
	defer func(testConnection net.Conn) {
		closureError := testConnection.Close()
		if closureError != nil {

		}
	}(testConnection)

}
