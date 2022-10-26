package main

import (
	"bufio"
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

// ChatUser is a basic "constructor" for each user.
type ChatUser struct {
	userConnection   net.Conn // userConnection is the socket.
	isDisconnected   bool
	userName         string
	outgoingMessages chan string   // Channel with all pending outgoing messages.
	ioReader         *bufio.Reader // Buffered socket read stream.
	ioWriter         *bufio.Writer // Buffered socket write stream.
}

func newChatUser(netConnection net.Conn) *ChatUser {
	writer := bufio.NewWriter(netConnection)
	reader := bufio.NewReader(netConnection)

	newUser := &ChatUser{
		userConnection:   netConnection,
		isDisconnected:   false,
		ioReader:         reader,
		ioWriter:         writer,
		outgoingMessages: make(chan string),
	}

	return newUser
}

func (chatRoom *ChatRoom) listenForMessages() {}
func (chatRoom *ChatRoom) Join(newUserConnection net.Conn) {
	newUser := newChatUser(newUserConnection)

	if newUser.Login(chatRoom) != nil {
		chatRoom.newJoins <- newUser
	}
}

func (chatUser *ChatUser) Login(chatRoom *ChatRoom) error {
	writeError := chatUser.WriteString("\nWelcome to the Couchbase Chat Server." +
		"\nAll the messages are end-to-end unencrypted." +
		"\nThe server is not sharing your private messages for commercial purposes, but only for I do not know how to...\n" +
		"Otherwise, I would definitely sell it to LEGO.\n\n")
	if writeError != nil {
		return writeError
	}

	writeError = chatUser.WriteString("Please enter your name:")
	if writeError != nil {
		return writeError
	}

	var userNameError error
	chatUser.userName, userNameError = chatUser.ReadLine()

	if userNameError != nil {
		return userNameError
	}

	log.Println("A New User Logged In:", chatUser.userName)
	writeError = chatUser.WriteString("Welcome " + chatUser.userName + "!")
	if writeError != nil {
		return writeError
	}

	return nil
}

func (chatUser *ChatUser) WriteString(messageToWrite string) error {
	_, writeError := chatUser.ioWriter.WriteString(messageToWrite)

	if writeError != nil {
		return writeError
	}

	return chatUser.ioWriter.Flush()
}

func (chatUser *ChatUser) ReadLine() (string, error) {
	inputBytes, _, inputError := chatUser.ioReader.ReadLine()

	return string(inputBytes), inputError
}

// Function main creates the socket and bind to port 8080 and wait for incoming connections using the loop.
// Everytime it receives a new connection, it passes it to" chatRoom.Join()".
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
