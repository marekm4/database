package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{}

func CreateDatabaseServer(database Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			query, err := ParseQuery(string(message))
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
			err = c.WriteMessage(mt, []byte(output))
			if err != nil {
				log.Println(err)
				break
			}
		}
	}
}

func main() {
	database := Database{make(map[string]any)}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/database", CreateDatabaseServer(database))
	port := "8080"
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
