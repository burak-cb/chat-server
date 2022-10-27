package main

import (
	"net"
	"testing"
	"time"
)

var chatRoom ChatRoom
var chatUser ChatUser

func init() {
	go StartServer()
	time.Sleep(time.Second)
}

func TestTCPServerRunning(t *testing.T) {
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
