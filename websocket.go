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

func ReloadDatabase(database Database) error {
	_, err := Exec(os.Getenv("DOWNLOAD_COMMAND"))
	if err != nil {
		return err
	}
	err = Load(database, "database.txt")
	if err != nil {
		return err
	}
	return nil
}

func ReloadRemoteDatabase() error {
	url := os.Getenv("RELOAD_URL")
	if len(url) > 0 {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		client := &http.Client{}
		_, err = client.Do(request)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	database := NewDatabase()
	err := ReloadDatabase(database)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		err := ReloadDatabase(database)
		if err != nil {
			log.Fatalln(err)
		}
	})
	http.HandleFunc("/database", DatabaseHandleFunc(database))

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func(database Database) {
		<-signals
		err := Dump(database, "database.txt")
		if err != nil {
			log.Fatalln(err)
		}
		_, err = Exec(os.Getenv("UPLOAD_COMMAND"))
		if err != nil {
			log.Fatalln(err)
		}
		err = ReloadRemoteDatabase()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}(database)

	port := "8080"
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
