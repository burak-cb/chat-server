package main

import (
	"ChatServer/experimentalChatServer/engineSide/handlerSide"
	"github.com/bmizerany/pat"
	"net/http"
)

func routes() http.Handler {
	multiPlexer := pat.New()
	multiPlexer.Get("/", http.HandlerFunc(handlerSide.HomePage))
	multiPlexer.Get("/ws", http.HandlerFunc(handlerSide.WebSocketEndpoint))

	return multiPlexer
}
