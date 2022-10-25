package main

import (
	"log"
	"net"
)

type ChatRoom struct {
	mapOfUsers       map[string]*ChatUser
	incomingMessages chan string
	newJoins         chan *ChatUser
	userDisconnects  chan string
}

type ChatUser struct{}

func newChatRoom() *ChatRoom {
	return &ChatRoom{
		mapOfUsers:       make(map[string]*ChatUser),
		incomingMessages: make(chan string),
		newJoins:         make(chan *ChatUser),
		userDisconnects:  make(chan string),
	}
}

func (chatRoom *ChatRoom) listenForMessages() {}
func (chatRoom *ChatRoom) Join(net.Conn)      {}

func main() {
	connectionServer, connectionError := net.Listen("tcp", "localhost:8080")
	if connectionError != nil {
		log.Println("Error Starting the Server:", connectionError)
	}

	log.Println("Spinning the Chat Server!")

	chatRoom := newChatRoom()
	chatRoom.listenForMessages()

	for {
		incomingConnection, newConnectionError := connectionServer.Accept()

		if newConnectionError != nil {
			log.Println("New Connection Error:", newConnectionError)
		} else {
			log.Println("New Connection:", incomingConnection.RemoteAddr().String())
		}

		chatRoom.Join(incomingConnection)
	}
}
