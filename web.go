package main

import (
	"bytes"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
)

func CreateDatabaseServer(database Database) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		for {
			buffer := make([]byte, 1024)
			_, err := ws.Read(buffer)
			if err != nil {
				log.Println(err)
			}
			query, err := ParseQuery(string(bytes.Trim(buffer, string([]byte{0}))))
			output := ""
			if err != nil {
				output = err.Error()
			} else {
				response, err := query.Execute(database)
				if err != nil {
					output = err.Error()
				} else {
					output = response
				}
			}
			_, err = ws.Write([]byte(output))
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func main() {
	database := Database{make(map[string]any)}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.Handle("/database", websocket.Handler(CreateDatabaseServer(database)))
	port := "8080"
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
