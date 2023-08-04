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

func (q EmptyQuery) Execute(database Database) []string {
	return []string{}
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
	Value int
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
		numericValue, _ := strconv.Atoi(value)
		return IncrementQuery{key, numericValue}
	}
	if operation == "append" {
		return AppendQuery{key, value}
	}
	return EmptyQuery{}
}
