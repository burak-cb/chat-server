package main

import (
	"bufio"
	"log"
	"net"
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

			writeError := chatUser.WriteString(messageData)
			if writeError != nil {
				log.Println("ERROR:", writeError)
			}
		}
	}()
}

func (chatUser *ChatUser) Login(chatRoom *ChatRoom) error {
	writeError := chatUser.WriteString("\nWelcome to the Couchbase Chat Server." +
		"\nAll the messages are end-to-end unencrypted." +
		"\nThe server is not sharing your private messages for commercial purposes, but only for I do not know how to...\n" +
		"Otherwise, I would definitely sell it to LEGO.\n\n")
	if writeError != nil {
		return writeError
	}

	writeError = chatUser.WriteString("Please enter your name: ")
	if writeError != nil {
		return writeError
	}

	var userNameError error
	chatUser.userName, userNameError = chatUser.ReadLine()

	if userNameError != nil {
		return userNameError
	}

	log.Println("A New User Logged In:", chatUser.userName)
	writeError = chatUser.WriteString("Welcome " + chatUser.userName + "!\n")
	if writeError != nil {
		return writeError
	}

	chatUser.WriteOutgoingMessages()
	chatUser.ReadIncomingMessages(chatRoom)
	return nil
}

func (chatUser *ChatUser) ReadLine() (string, error) {
	inputBytes, _, inputError := chatUser.ioReader.ReadLine()
	inputString := string(inputBytes)
	return inputString, inputError
}

func (chatUser *ChatUser) WriteString(messageToWrite string) error {
	_, writeError := chatUser.ioWriter.WriteString(messageToWrite)

	if writeError != nil {
		return writeError
	}

	return chatUser.ioWriter.Flush()
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
}
