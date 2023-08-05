package main

import (
	"github.com/gorilla/websocket"
	"gotest.tools/v3/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebsocket(t *testing.T) {
	// Given empty database
	database := NewDatabase()
	server := httptest.NewServer(http.HandlerFunc(DatabaseHandleFunc(database)))
	ws, _, err := websocket.DefaultDialer.Dial(strings.Replace(server.URL, "http", "ws", 1), nil)
	if err != nil {
		assert.NilError(t, err)
	}

	// When we add value
	err = ws.WriteMessage(websocket.TextMessage, []byte("update username john"))
	if err != nil {
		assert.NilError(t, err)
	}

	// Then value is there
	err = ws.WriteMessage(websocket.TextMessage, []byte("select username"))
	if err != nil {
		assert.NilError(t, err)
	}
	_, message, err := ws.ReadMessage()
	if err != nil {
		assert.NilError(t, err)
	}
	assert.Equal(t, string(message), "john")

	// When we increment value
	err = ws.WriteMessage(websocket.TextMessage, []byte("increment money 100"))
	if err != nil {
		assert.NilError(t, err)
	}

	// Then value is there
	err = ws.WriteMessage(websocket.TextMessage, []byte("select money"))
	if err != nil {
		assert.NilError(t, err)
	}
	_, message, err = ws.ReadMessage()
	if err != nil {
		assert.NilError(t, err)
	}
	assert.Equal(t, string(message), "100.000000")

	// When we add value
	err = ws.WriteMessage(websocket.TextMessage, []byte("append orders pizza"))
	if err != nil {
		assert.NilError(t, err)
	}

	// Then value is there
	err = ws.WriteMessage(websocket.TextMessage, []byte("select orders"))
	if err != nil {
		assert.NilError(t, err)
	}
	_, message, err = ws.ReadMessage()
	if err != nil {
		assert.NilError(t, err)
	}
	assert.Equal(t, string(message), "pizza")
}
