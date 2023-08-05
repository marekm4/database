package main

import (
	"github.com/gorilla/websocket"
	"gotest.tools/v3/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestWebsocket(t *testing.T) {
	// Given empty database server
	_, mux, err := NewContainer("test.txt")
	assert.NilError(t, err)
	server := httptest.NewServer(mux)
	ws, _, err := websocket.DefaultDialer.Dial(strings.Replace(server.URL, "http", "ws", 1)+"/database", nil)
	assert.NilError(t, err)

	// When we add value
	err = ws.WriteMessage(websocket.TextMessage, []byte("update username john"))
	assert.NilError(t, err)

	// Then value is there
	err = ws.WriteMessage(websocket.TextMessage, []byte("select username"))
	assert.NilError(t, err)
	_, message, err := ws.ReadMessage()
	assert.NilError(t, err)
	assert.Equal(t, string(message), "john")

	// When we increment value
	err = ws.WriteMessage(websocket.TextMessage, []byte("increment money 100"))
	assert.NilError(t, err)

	// Then value is there
	err = ws.WriteMessage(websocket.TextMessage, []byte("select money"))
	assert.NilError(t, err)
	_, message, err = ws.ReadMessage()
	assert.NilError(t, err)
	assert.Equal(t, string(message), "100.000000")

	// When we add value
	err = ws.WriteMessage(websocket.TextMessage, []byte("append orders pizza"))
	assert.NilError(t, err)

	// Then value is there
	err = ws.WriteMessage(websocket.TextMessage, []byte("select orders"))
	assert.NilError(t, err)
	_, message, err = ws.ReadMessage()
	assert.NilError(t, err)
	assert.Equal(t, string(message), "pizza")
}

func TestWebsocket_Reload(t *testing.T) {
	// Given empty database server
	_, mux, err := NewContainer("test.txt")
	assert.NilError(t, err)
	server := httptest.NewServer(mux)
	ws, _, err := websocket.DefaultDialer.Dial(strings.Replace(server.URL, "http", "ws", 1)+"/database", nil)
	assert.NilError(t, err)

	// When we add value
	err = ws.WriteMessage(websocket.TextMessage, []byte("update username john"))
	assert.NilError(t, err)

	// Then value is there
	err = ws.WriteMessage(websocket.TextMessage, []byte("select username"))
	assert.NilError(t, err)
	_, message, err := ws.ReadMessage()
	assert.NilError(t, err)
	assert.Equal(t, string(message), "john")

	// When we reload server
	err = ReloadRemoteDatabase(server.URL + "/reload")
	assert.NilError(t, err)

	// Then value is gone
	err = ws.WriteMessage(websocket.TextMessage, []byte("select username"))
	assert.NilError(t, err)
	_, message, err = ws.ReadMessage()
	assert.NilError(t, err)
	assert.Equal(t, string(message), "")
}

func TestWebsocket_Client(t *testing.T) {
	// Given empty database server
	_, mux, err := NewContainer("test.txt")
	assert.NilError(t, err)
	server := httptest.NewServer(mux)

	// When we ask about client
	request, err := http.NewRequest("GET", server.URL, nil)
	assert.NilError(t, err)
	client := &http.Client{}
	response, err := client.Do(request)
	assert.NilError(t, err)

	// Then client is there
	body, err := io.ReadAll(response.Body)
	assert.NilError(t, err)
	assert.Check(t, strings.Contains(string(body), "<title>Database</title>"))
}

func TestWebsocket_Shutdown(t *testing.T) {
	// Given database with data
	database := NewDatabase()
	database.Update("username", "john")
	filename := "test.txt"

	// When we shut it down
	err := Shutdown(database, filename)
	assert.NilError(t, err)

	// Then state is saved
	database = NewDatabase()
	err = Load(database, filename)
	assert.NilError(t, err)
	assert.DeepEqual(t, database.Select("username"), []string{"john"})

	// Clean up
	err = os.Remove(filename)
	assert.NilError(t, err)
}
