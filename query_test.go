package main

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestParseQuery_Ok(t *testing.T) {
	testCases := map[string]struct {
		Query    string
		Expected Query
	}{
		"Empty":     {"select", EmptyQuery{}},
		"Select":    {"select user_name", SelectQuery{"user_name"}},
		"Update":    {"update user_name john", UpdateQuery{"user_name", "john"}},
		"Increment": {"increment user_money 5", IncrementQuery{"user_money", 5}},
		"Append":    {"append user_history order 1", AppendQuery{"user_history", "order 1"}},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			query := ParseQuery(testCase.Query)
			assert.DeepEqual(t, query, testCase.Expected)
		})
	}
}
