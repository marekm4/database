package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Database struct {
	Data map[string]any
}

type Query interface {
	Execute(Database) (string, error)
}

type GetQuery struct {
	Key string
}

func (q GetQuery) Execute(database Database) (string, error) {
	return "getting", nil
}

type SetQuery struct {
	Key   string
	Value string
}

func (q SetQuery) Execute(database Database) (string, error) {
	return "setting", nil
}

type IncrementQuery struct {
	Key   string
	Value int
}

func (q IncrementQuery) Execute(database Database) (string, error) {
	return "incrementing", nil
}

type AppendQuery struct {
	Key   string
	Value string
}

func (q AppendQuery) Execute(database Database) (string, error) {
	return "appending", nil
}

func ParseQuery(query string) (Query, error) {
	i := strings.Index(query, " ")
	if i < 0 {
		return nil, errors.New("invalid query")
	}
	operation := query[:i]
	key := query[i+1:]
	if operation == "get" {
		return GetQuery{key}, nil
	}
	i = strings.Index(key, " ")
	if i < 0 {
		return nil, errors.New("not enough arguments")
	}
	value := key[i+1:]
	key = key[:i]
	if operation == "set" {
		return SetQuery{key, value}, nil
	}
	if operation == "increment" {
		fmt.Println(value)
		numericValue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println(err.Error())
			return nil, errors.New("invalid numeric value")
		}
		return IncrementQuery{key, numericValue}, nil
	}
	if operation == "append" {
		return AppendQuery{key, value}, nil
	}
	return nil, errors.New("invalid operation")
}
