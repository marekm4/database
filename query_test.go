package main

import "testing"

func TestParseQuery_Ok(t *testing.T) {
	testCases := map[string]struct {
		Query    string
		Expected Query
	}{
		"Empty":                {"get", EmptyQuery{}},
		"Get":                  {"get user_name", GetQuery{"user_name"}},
		"Set":                  {"set user_name john", SetQuery{"user_name", "john"}},
		"Increment":            {"increment user_money 5", IncrementQuery{"user_money", 5}},
		"Increment_non_nmeric": {"increment user_money a", IncrementQuery{"user_money", 0}},
		"Append":               {"append user_history order 1", AppendQuery{"user_history", "order 1"}},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			query := ParseQuery(testCase.Query)
			if query != testCase.Expected {
				t.Fatalf("%v expected, got %v", testCase.Expected, query)
			}
		})
	}
}
