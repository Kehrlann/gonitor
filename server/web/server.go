package web

import (
	"net/http"
	"github.com/kehrlann/gonitor/server/web/handlers"
	"github.com/gorilla/websocket"
)

func Serve() <-chan *websocket.Conn {
	incoming_connections := make(chan *websocket.Conn)
	websocket_hander := handlers.WebsocketHandler{incoming_connections}
	serve(websocket_hander)
	return incoming_connections
}

func serve(websocketHandler http.Handler) func() {
	http.HandleFunc("/", handlers.HandleIndex)
	http.Handle("/ws", websocketHandler)
	server := &http.Server{Addr:":3000",Handler:nil}
	go server.ListenAndServe()

	shutDown := func () { server.Close() }
	return shutDown
}
