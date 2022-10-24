module ChatServer

go 1.19

require (
	// HTML Template Package.
	github.com/CloudyKit/fastprinter v0.0.0-20200109182630-33d98a066a53 // indirect
	github.com/CloudyKit/jet/v6 v6.1.0 // indirect
)

// indirect - Routing package for directing traffic from web-server to the appropriate handler.
require github.com/bmizerany/pat v0.0.0-20210406213842-e4b6760bdd6f

// Actual WebSocket Package.
require github.com/gorilla/websocket v1.5.0 // indirect