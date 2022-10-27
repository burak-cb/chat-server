package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"time"
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
				splitString := strings.Split(incomingMessage, " ")
				if splitString[2] == "/startStory" {
					chatRoom.BroadcastMessage("Douglas Noel Adams is here to tell you an amazing story.\n")
					go chatRoom.DNAStory(splitString[0])
				} else if splitString[2] == "/whisper" {
					if chatRoom.mapOfUsers[splitString[3]] != nil {
						whisperError := chatRoom.mapOfUsers[splitString[3]].
							WriteString(splitString[0] + " (Whispers) -> " + strings.Join(splitString[4:], " "))
						if whisperError != nil {
							log.Println("ERROR:", whisperError)
						}
					} else {
						whisperError := chatRoom.mapOfUsers[splitString[0]].WriteString(splitString[3] + " is not online.")
						if whisperError != nil {
							log.Println("ERROR:", whisperError)
						}
					}
				} else {
					chatRoom.BroadcastMessage(incomingMessage)
				}

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

func (chatRoom *ChatRoom) DNAStory(askingUser string) {
	inputFile, fileOpeningError := os.Open("hitchhikersGuide.txt")

	if fileOpeningError != nil {
		log.Fatalln("ERROR:", fileOpeningError)
	}

	defer func(inputFile *os.File) {
		fileClosureError := inputFile.Close()
		if fileClosureError != nil {

		}
	}(inputFile)

	fileScanner := bufio.NewScanner(inputFile)

	for fileScanner.Scan() && chatRoom.mapOfUsers[askingUser] != nil {
		chatRoom.BroadcastMessage(fileScanner.Text())
		time.Sleep(time.Second)
	}

	chatRoom.BroadcastMessage("\nDNA is stopping because " + askingUser + " has left the server.\n")
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
