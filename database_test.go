package main

import (
	"fmt"
	"gotest.tools/v3/assert"
	"testing"
)

func TestDatabase_GetEmpty(t *testing.T) {
	// Given empty database
	database := NewDatabase()

	// When we ask for not existing key
	values := database.Select("not_exists")

	// Then key does not exist
	assert.DeepEqual(t, values, []string{""})
}

func TestDatabase_Update(t *testing.T) {
	// Given empty database
	database := NewDatabase()
	username := "john"

	// When we update value
	database.Update("username", username)

	// Then value is set
	values := database.Select("username")
	assert.DeepEqual(t, values, []string{username})

	// When we update value
	newUsername := "alice"
	database.Update("username", newUsername)

	// Then value is updated
	values = database.Select("username")
	assert.DeepEqual(t, values, []string{newUsername})
}

func TestDatabase_Increment(t *testing.T) {
	// Given empty database
	database := NewDatabase()
	money := 5.5

	// When we increment counter
	database.Increment("money", money)

	// Then counter is set
	values := database.Select("money")
	assert.DeepEqual(t, values, []string{fmt.Sprintf("%f", money)})

	// When we increment counter
	database.Increment("money", money)

	// Then counter is incremented
	values = database.Select("money")
	assert.DeepEqual(t, values, []string{fmt.Sprintf("%f", money*2)})
}

func TestDatabase_Append(t *testing.T) {
	// Given empty database
	database := NewDatabase()
	firstOrder := "order 1"

	// When we append value
	database.Append("orders", firstOrder)

	// Then value is appended
	values := database.Select("orders")
	assert.DeepEqual(t, values, []string{firstOrder})

	// When we append another value
	secondOrder := "order 2"
	database.Append("orders", secondOrder)

	// Then both values are there
	values = database.Select("orders")
	assert.DeepEqual(t, values, []string{firstOrder, secondOrder})
}

func TestDatabase_Clear(t *testing.T) {
	// Given database with data
	database := NewDatabase()
	key := "username"
	database.Update(key, "john")

	// When database is cleared
	database.Clear()

	// Then data is gone
	values := database.Select(key)
	assert.DeepEqual(t, values, []string{""})
}
