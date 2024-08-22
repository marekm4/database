package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	delim := byte('\n')
	database := NewDatabase()

	filename := "database.txt"
	if len(os.Getenv("FILE")) > 0 {
		filename = os.Getenv("FILE")
	}

	err := Load(database, filename)
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

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println(err)
			break
		}

		go func(net.Conn, Database) {
			reader := bufio.NewReader(connection)
			writer := bufio.NewWriter(connection)

			for {
				message, err := reader.ReadString(delim)
				if err != nil {
					log.Println(err)
					break
				}

				query := ParseQuery(strings.TrimSuffix(message, string(delim)))
				values := query.Execute(database)

				_, err = writer.WriteString(strings.Join(values, "\n") + string(delim))
				if err != nil {
					log.Println(err)
					break
				}
				err = writer.Flush()
				if err != nil {
					log.Println(err)
					break
				}
			}
			err = connection.Close()
			if err != nil {
				log.Println(err)
			}
		}(connection, database)
	}

	err = listener.Close()
	if err != nil {
		log.Println(err)
	}
}
