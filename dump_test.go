package main

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestDumpQueries(t *testing.T) {
	// Given database with values
	database := NewDatabase()
	database.Update("username", "john")
	database.Increment("age", 30)
	database.Increment("money", 100)
	database.Append("orders", "pizza")
	database.Append("orders", "burgers")

	// When we dump it
	queries := DumpQueries(database)

	// Then we have queries
	assert.DeepEqual(t, queries, []string{
		"update username john",
		"increment age 30",
		"increment money 100",
		"append orders pizza",
		"append orders burgers",
	})
}

func TestLoadQueries(t *testing.T) {
	// Given empty database
	database := NewDatabase()
	queries := []string{
		"update username john",
		"increment age 30",
		"increment money 100",
		"append orders pizza",
		"append orders burgers",
	}

	// When we load queries
	LoadQueries(database, queries)

	// Then we have them in database
	assert.DeepEqual(t, database.Select("username"), []string{"john"})
	assert.DeepEqual(t, database.Select("age"), []string{"30"})
	assert.DeepEqual(t, database.Select("money"), []string{"100"})
	assert.DeepEqual(t, database.Select("orders"), []string{"pizza", "burgers"})
}
