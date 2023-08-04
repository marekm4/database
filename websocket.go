package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

var Upgrader = websocket.Upgrader{}

func DatabaseHandleFunc(database Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		connection, err := Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		for {
			messageType, message, err := connection.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			query := ParseQuery(string(message))
			log.Println(query)
			output := ""
			response, err := query.Execute(database)
			if err != nil {
				output = err.Error()
			} else {
				output = response
			}
			err = connection.WriteMessage(messageType, []byte(output))
			if err != nil {
				log.Println(err)
				break
			}
		}
		err = connection.Close()
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	database := NewDatabase()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/database", DatabaseHandleFunc(database))
	port := "8080"
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
