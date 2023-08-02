package main

import (
	"io"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

func DatabaseServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.Handle("/database", websocket.Handler(DatabaseServer))
	port := "8080"
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}
	http.ListenAndServe(":"+port, nil)
}
