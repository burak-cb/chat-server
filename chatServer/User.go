package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
)

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

func (chatUser *ChatUser) ReadIncomingMessages(chatRoom *ChatRoom) {
	go func() {
		for {
			incomingLine, messageReadError := chatUser.ReadLine()

			if chatUser.isDisconnected {
				break
			}

			if messageReadError != nil {
				chatRoom.LogOut(chatUser.userName)
				break
			}

			if incomingLine != "" {
				chatRoom.incomingMessages <- chatUser.userName + " -> " + incomingLine
			}
		}
	}()
}

func (chatUser *ChatUser) WriteOutgoingMessages() {
	go func() {
		for {
			messageData := <-chatUser.outgoingMessages
			messageData = messageData + "\n"

			chatUser.WriteString(messageData)
		}
	}()
}

func (chatUser *ChatUser) Login(chatRoom *ChatRoom) error {
	chatUser.WriteString("\nWelcome to the Couchbase Chat Server." +
		"\nAll the messages are end-to-end unencrypted." +
		"\nThe server is not sharing your private messages for commercial purposes, but only for I do not know how to...\n" +
		"Otherwise, I would definitely sell it to LEGO.\n\n")

	chatUser.WriteString("Please enter your name: ")

	var userNameError error
	chatUser.userName, userNameError = chatUser.ReadLine()

	if userNameError != nil {
		return userNameError
	}

	serverPassword := "dontPanic!"

	for e := 2; e >= 0; e-- {
		chatUser.WriteString("Server Password: ")

		passwordInput, passwordReadingError := chatUser.ReadLine()

		if passwordReadingError != nil {
			return passwordReadingError
		}

		if passwordInput == serverPassword {
			break
		} else if passwordInput != serverPassword {
			chatUser.WriteString("The password is incorrect.\n")
		}

		if e == 0 {
			chatUser.WriteString("Access is not granted.\n")
			chatUser.Close()
		}

		chatUser.WriteString("You have " + strconv.Itoa(e) + " more attempt.\n\n")
	}

	chatUser.WriteString("Access granted.\n\n")

	log.Println("A New User Logged In:", chatUser.userName)
	chatUser.WriteString("Welcome " + chatUser.userName + "!\n")

	chatUser.WriteOutgoingMessages()
	chatUser.ReadIncomingMessages(chatRoom)
	return nil
}

func (chatUser *ChatUser) ReadLine() (string, error) {
	inputBytes, _, inputError := chatUser.ioReader.ReadLine()
	inputString := string(inputBytes)
	return inputString, inputError
}

func (chatUser *ChatUser) WriteString(messageToWrite string) {
	// _, writeError := chatUser.ioWriter.WriteString("\0337\033[A\033[999D\033[S\033[L" + messageToWrite + "\0338")
	_, writeError := chatUser.ioWriter.WriteString(messageToWrite)

	if writeError != nil {
		chatUser.WriteString("The server has failed to send the message.")
		log.Println("Message Write Error:", writeError)
	}

	messageFlushError := chatUser.ioWriter.Flush()
	if messageFlushError != nil {
		log.Println("Message Flush Error:", writeError)
	}
}

func (chatUser *ChatUser) SendMessage(messageToSend string) {
	chatUser.outgoingMessages <- messageToSend
}

func (chatUser *ChatUser) Close() {
	chatUser.isDisconnected = true
	disconnectionError := chatUser.userConnection.Close()
	if disconnectionError != nil {
		return
	}

	log.Println(chatUser.userName, "has disconnected.")
}
