package main

import (
	"errors"
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
