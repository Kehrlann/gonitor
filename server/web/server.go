package web

import (
	"net/http"
	"github.com/kehrlann/gonitor/server/web/assets"
)

func Serve() func() {
	// TODO : how to bundle these assets ?
	//fs := http.FileServer(http.Dir("server/web/assets"))
	http.HandleFunc("/", assets.HandleIndex)
	server := &http.Server{Addr:":3000",Handler:nil}
	go server.ListenAndServe()

	shutDown := func () { server.Close() }
	return shutDown
}
