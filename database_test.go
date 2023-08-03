package main

import (
	"errors"
	"strconv"
	"strings"
	"testing"
)

func TestParseQuery_Ok(t *testing.T) {
	testCases := map[string]struct {
		Query    string
		Expected Query
	}{
		"get":       {"get user_name", GetQuery{"user_name"}},
		"set":       {"set user_name john", SetQuery{"user_name", "john"}},
		"increment": {"increment user_money 5", IncrementQuery{"user_money", 5}},
		"append":    {"append user_history order 1", AppendQuery{"user_history", "order 1"}},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			query, err := ParseQuery(testCase.Query)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if query != testCase.Expected {
				t.Fatalf("%v expected, got %v", testCase.Expected, query)
			}
		})
	}
}

func TestParseQuery_Err(t *testing.T) {
	testCases := map[string]struct {
		Query       string
		ExpectedErr error
	}{
		"empty query":           {"", errors.New("invalid query")},
		"not enough arguments":  {"set user_name", errors.New("not enough arguments")},
		"invalid operation":     {"delete user_name 5", errors.New("invalid operation")},
		"invalid numeric value": {"increment user_money string", errors.New("invalid numeric value")},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := ParseQuery(testCase.Query)
			if testCase.ExpectedErr.Error() != err.Error() {
				t.Fatalf("%v expected, got %v", testCase.ExpectedErr, err)
			}
		})
	}
}

func TestDatabase_GetQuery_NotExist(t *testing.T) {
	// Given empty database
	database := NewDatabase()

	// When we ask for not existing key
	_, err := GetQuery{"user_name"}.Execute(database)

	// Then key does not exist
	expectedErr := errors.New("key not exists")
	if expectedErr.Error() != err.Error() {
		t.Fatalf("%v expected, got %v", expectedErr, err)
	}
}

func TestDatabase_SetQuery(t *testing.T) {
	// Given empty database
	database := NewDatabase()
	userName := "john"

	// When we set value
	_, err := SetQuery{"user_name", userName}.Execute(database)
	if err != nil {
		t.Fatalf("got unexpected %v", err)
	}

	// Then value is set
	result, err := GetQuery{"user_name"}.Execute(database)
	if err != nil {
		t.Fatalf("got unexpected %v", err)
	}
	if userName != result {
		t.Fatalf("%v expected, got %v", userName, result)
	}

	// When we update value
	newUserName := "alice"
	_, err = SetQuery{"user_name", newUserName}.Execute(database)
	if err != nil {
		t.Fatalf("got unexpected %v", err)
	}

	// Then value is updated
	result, err = GetQuery{"user_name"}.Execute(database)
	if err != nil {
		t.Fatalf("got unexpected %v", err)
	}
	if newUserName != result {
		t.Fatalf("%v expected, got %v", newUserName, result)
	}
}

func TestDatabase_SetQuery_OverCounter(t *testing.T) {
	// Given counter
	database := NewDatabase()
	_, err := IncrementQuery{"user_money", 5}.Execute(database)
	if err != nil {
		t.Fatalf("got unexpected %v", err)
	}

	// When we try to set string over counter
	_, err = SetQuery{"user_money", "john"}.Execute(database)

	// Then it is wrong type
	expectedErr := errors.New("wrong type")
	if expectedErr.Error() != err.Error() {
		t.Fatalf("%v expected, got %v", expectedErr, err)
	}
}

func TestDatabase_IncrementQuery(t *testing.T) {
	// Given empty database
	database := NewDatabase()
	money := 5

	// When we set counter
	_, err := IncrementQuery{"user_money", money}.Execute(database)
	if err != nil {
		t.Fatalf("got unexpected %v", err)
	}

	// Then counter is set
	result, err := GetQuery{"user_money"}.Execute(database)
	if err != nil {
		t.Fatalf("got unexpected %v", err)
	}
	if strconv.Itoa(money) != result {
		t.Fatalf("%v expected, got %v", strconv.Itoa(money), result)
	}

	// When we increment counter
	_, err = IncrementQuery{"user_money", money}.Execute(database)
	if err != nil {
		t.Fatalf("got unexpected %v", err)
	}

	// Then counter is incremented
	result, err = GetQuery{"user_money"}.Execute(database)
	if err != nil {
		t.Fatalf("got unexpected %v", err)
	}
	if strconv.Itoa(money*2) != result {
		t.Fatalf("%v expected, got %v", strconv.Itoa(money*2), result)
	}
}

func TestDatabase_IncrementQuery_OverString(t *testing.T) {
	// Given string
	database := NewDatabase()
	_, err := SetQuery{"user_name", "john"}.Execute(database)
	if err != nil {
		t.Fatalf("got unexpected %v", err)
	}

	// When we try to increment counter over string
	_, err = IncrementQuery{"user_name", 5}.Execute(database)

	// Then it is wrong type
	expectedErr := errors.New("wrong type")
	if expectedErr.Error() != err.Error() {
		t.Fatalf("%v expected, got %v", expectedErr, err)
	}
}

func TestDatabase_AppendQuery(t *testing.T) {
	// Given empty database
	database := NewDatabase()
	firstOrder := "order 1"

	// When we append value
	_, err := AppendQuery{"user_orders", firstOrder}.Execute(database)
	if err != nil {
		t.Fatalf("got unexpected %v", err)
	}

	// Then value is there
	result, err := GetQuery{"user_orders"}.Execute(database)
	if err != nil {
		t.Fatalf("got unexpected %v", err)
	}
	if firstOrder != result {
		t.Fatalf("%v expected, got %v", firstOrder, result)
	}

	// When we append another string
	secondOrder := "order 2"
	_, err = AppendQuery{"user_orders", secondOrder}.Execute(database)
	if err != nil {
		t.Fatalf("got unexpected %v", err)
	}

	// Then both strings are there
	result, err = GetQuery{"user_orders"}.Execute(database)
	if err != nil {
		t.Fatalf("got unexpected %v", err)
	}
	if strings.Join([]string{firstOrder, secondOrder}, "\n") != result {
		t.Fatalf("%v expected, got %v", strings.Join([]string{firstOrder, secondOrder}, "\n"), result)
	}
}
