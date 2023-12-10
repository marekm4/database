package main

import (
	"strconv"
	"strings"
)

type Query interface {
	Execute(Database) []string
}

type EmptyQuery struct {
}

func (q EmptyQuery) Execute(_ Database) []string {
	return []string{}
}

type ListQuery struct {
	Key string
}

func (q ListQuery) Execute(database Database) []string {
	return database.List(q.Key)
}

type SelectQuery struct {
	Key string
}

func (q SelectQuery) Execute(database Database) []string {
	return database.Select(q.Key)
}

type UpdateQuery struct {
	Key   string
	Value string
}

func (q UpdateQuery) Execute(database Database) []string {
	database.Update(q.Key, q.Value)
	return []string{}
}

type IncrementQuery struct {
	Key   string
	Value float64
}

func (q IncrementQuery) Execute(database Database) []string {
	database.Increment(q.Key, q.Value)
	return []string{}
}

type AppendQuery struct {
	Key   string
	Value string
}

func (q AppendQuery) Execute(database Database) []string {
	database.Append(q.Key, q.Value)
	return []string{}
}

func ParseQuery(query string) Query {
	i := strings.Index(query, " ")
	if i < 0 {
		return EmptyQuery{}
	}
	operation := query[:i]
	key := query[i+1:]
	if operation == "list" {
		return ListQuery{key}
	}
	if operation == "select" {
		return SelectQuery{key}
	}
	i = strings.Index(key, " ")
	if i < 0 {
		return EmptyQuery{}
	}
	value := key[i+1:]
	key = key[:i]
	if operation == "update" {
		return UpdateQuery{key, value}
	}
	if operation == "increment" {
		numericValue, _ := strconv.ParseFloat(value, 64)
		return IncrementQuery{key, numericValue}
	}
	if operation == "append" {
		return AppendQuery{key, value}
	}
	return EmptyQuery{}
}
