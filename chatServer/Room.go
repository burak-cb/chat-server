package main

import (
	"log"
	"net"
)

type ChatRoom struct { // Main room of the server.
	mapOfUsers       map[string]*ChatUser // Stores all the users and their connection information.
	incomingMessages chan string          // Receives incoming messages from ChatUser connection.
	newJoins         chan *ChatUser       // Receives incoming new ChatUser connections.
	userDisconnects  chan string          // Receives disconnect notifications.
}

// newChatroom is a basic "constructor" for each user.
func newChatRoom() *ChatRoom {
	return &ChatRoom{
		mapOfUsers:       make(map[string]*ChatUser),
		incomingMessages: make(chan string),
		newJoins:         make(chan *ChatUser),
		userDisconnects:  make(chan string),
	}
}

func (chatRoom *ChatRoom) MessageListener() {
	go func() {
		for {
			select {
			case incomingMessage := <-chatRoom.incomingMessages:
				chatRoom.BroadcastMessage(incomingMessage)
			case joiningUser := <-chatRoom.newJoins:
				chatRoom.mapOfUsers[joiningUser.userName] = joiningUser
				chatRoom.BroadcastMessage(joiningUser.userName + " has just joined the chat room.")
			case disconnectedUser := <-chatRoom.userDisconnects:
				if chatRoom.mapOfUsers[disconnectedUser] != nil {
					chatRoom.mapOfUsers[disconnectedUser].Close()
					delete(chatRoom.mapOfUsers, disconnectedUser)
					chatRoom.BroadcastMessage(disconnectedUser + " has disconnected from the server.")
				}
			}
		}
	}()
}

func (chatRoom *ChatRoom) LogOut(disconnectedUser string) {
	chatRoom.userDisconnects <- disconnectedUser
}

func (chatRoom *ChatRoom) JoinServer(newUserConnection net.Conn) {
	newUser := newChatUser(newUserConnection)

	if newUser.Login(chatRoom) == nil {
		chatRoom.newJoins <- newUser
	}
}
func (chatRoom *ChatRoom) BroadcastMessage(messageToBroadcast string) {
	for _, eachUser := range chatRoom.mapOfUsers {
		eachUser.SendMessage(messageToBroadcast)
	}
}

func StartServer() {
	connectionServer, connectionError := net.Listen("tcp", "localhost:8080")
	if connectionError != nil {
		log.Println("Error Starting the Server:", connectionError)
	}

	log.Println("Spinning the Chat Server!")

	chatRoom := newChatRoom()
	chatRoom.MessageListener()

	for {
		incomingConnection, newConnectionError := connectionServer.Accept()

		if newConnectionError != nil {
			log.Println("New Connection Error:", newConnectionError)
		} else {
			log.Println("New Connection:", incomingConnection.RemoteAddr().String())
		}

		go chatRoom.JoinServer(incomingConnection) // Go routine ensures multiple connections even at the exact same moment.
	}
}
