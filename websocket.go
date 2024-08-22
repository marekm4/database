package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
			values := query.Execute(database)
			if len(values) > 0 {
				err = connection.WriteMessage(messageType, []byte(strings.Join(values, "\n")))
				if err != nil {
					log.Println(err)
					break
				}
			}
		}
		err = connection.Close()
		if err != nil {
			log.Println(err)
		}
	}
}

func NewContainer(filename string) (Database, *http.ServeMux, error) {
	database := NewDatabase()
	err := Load(database, filename)
	if err != nil {
		return database, nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	mux.HandleFunc("/database", DatabaseHandleFunc(database))

	return database, mux, nil
}

func main() {
	filename := "database.txt"
	database, mux, err := NewContainer(filename)
	if err != nil {
		log.Fatalln(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		err := Dump(database, filename)
		if err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)
	}()

	port := "8080"
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalln(err)
	}
}
