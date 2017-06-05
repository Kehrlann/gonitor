package web

import (
	"net/http"
	"github.com/kehrlann/gonitor/server/web/handlers"
)

func Serve() func() {
	http.HandleFunc("/", handlers.HandleIndex)
	http.HandleFunc("/ws", handlers.HandleWebsockets)
	server := &http.Server{Addr:":3000",Handler:nil}
	go server.ListenAndServe()

	shutDown := func () { server.Close() }
	return shutDown
}
