package handlerSide

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var jetView = jet.NewSet(
	jet.NewOSFileSystemLoader("./HTML"),
	jet.InDevelopmentMode(), // Removes the need to re-start the application. Take it out before the production.
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin:     func(httpRequest *http.Request) bool { return true },
}

func HomePage(responseWriter http.ResponseWriter, pointerToRequest *http.Request) {
	renderPageError := pageRenderer(responseWriter, "homePageWithTemplate.jet", nil)

	if renderPageError != nil {
		log.Println(renderPageError)
	}
}

// WsJsonResponse defines the response sent back from web-socket.
type wsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"messageType"`
}

// WebSocketEndpoint upgrades connection to web-socket.
func WebSocketEndpoint(responseWriter http.ResponseWriter, pointerToRequest *http.Request) {
	webSocketConnection, wsError := upgradeConnection.Upgrade(responseWriter, pointerToRequest, nil)
	if wsError != nil {
		log.Print(wsError)
	}

	log.Println("A Client is connected to the endpoint.")

	var wsResponse wsJsonResponse
	wsResponse.Message = `<em><small>Connection established</small></em>`

	wsError = webSocketConnection.WriteJSON(wsResponse)
	if wsError != nil {
		log.Print(wsError)
	}
}

func pageRenderer(responseWriter http.ResponseWriter, templateToRender string, templateData jet.VarMap) error {
	pageView, pageError := jetView.GetTemplate(templateToRender)

	if pageError != nil {
		log.Println(pageError)
		return pageError
	}

	pageError = pageView.Execute(responseWriter, templateData, nil)
	if pageError != nil {
		log.Println(pageError)
		return pageError
	}

	return nil
}
