package main

import (
	"gotest.tools/v3/assert"
	"os"
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
		"increment age 30.000000",
		"increment money 100.000000",
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
	assert.DeepEqual(t, database.Select("age"), []string{"30.000000"})
	assert.DeepEqual(t, database.Select("money"), []string{"100.000000"})
	assert.DeepEqual(t, database.Select("orders"), []string{"pizza", "burgers"})
}

func TestFiles(t *testing.T) {
	// Given database with values
	database := NewDatabase()
	database.Update("username", "john")
	database.Increment("age", 30)
	database.Increment("money", 100)
	database.Append("orders", "pizza")
	database.Append("orders", "burgers")

	// When we dump it to file
	filename := "test.txt"
	err := Dump(database, filename)
	assert.NilError(t, err)

	// Then file exist
	_, err = os.Stat(filename)
	assert.NilError(t, err)

	// When we load it to new database
	database = NewDatabase()
	err = Load(database, filename)
	assert.NilError(t, err)

	// Then records are there
	assert.DeepEqual(t, database.Select("username"), []string{"john"})
	assert.DeepEqual(t, database.Select("age"), []string{"30.000000"})
	assert.DeepEqual(t, database.Select("money"), []string{"100.000000"})
	assert.DeepEqual(t, database.Select("orders"), []string{"pizza", "burgers"})

	// Clean up
	err = os.Remove(filename)
	assert.NilError(t, err)
}
