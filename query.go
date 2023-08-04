package main

import (
	"errors"
	"strconv"
	"strings"
)

type Query interface {
	Execute(Database) (string, error)
}

type EmptyQuery struct {
}

func (q EmptyQuery) Execute(database Database) (string, error) {
	return "", nil
}

type GetQuery struct {
	Key string
}

func (q GetQuery) Execute(database Database) (string, error) {
	value, exists := database.Data[q.Key]
	if !exists {
		return "", errors.New("key not exists")
	}
	if exists {
		if parsed, ok := value.(string); ok {
			return parsed, nil
		}
		if parsed, ok := value.(int); ok {
			return strconv.Itoa(parsed), nil
		}
		if parsed, ok := value.([]string); ok {
			return strings.Join(parsed, "\n"), nil
		}
	}
	return "", nil
}

type SetQuery struct {
	Key   string
	Value string
}

func (q SetQuery) Execute(database Database) (string, error) {
	value, exists := database.Data[q.Key]
	if exists {
		if _, ok := value.(string); !ok {
			return "", errors.New("wrong type")
		}
	}
	database.Data[q.Key] = q.Value
	return "", nil
}

type IncrementQuery struct {
	Key   string
	Value int
}

func (q IncrementQuery) Execute(database Database) (string, error) {
	value, exists := database.Data[q.Key]
	if exists {
		parsed, ok := value.(int)
		if !ok {
			return "", errors.New("wrong type")
		} else {
			database.Data[q.Key] = parsed + q.Value
		}
	} else {
		database.Data[q.Key] = q.Value
	}
	return "", nil
}

type AppendQuery struct {
	Key   string
	Value string
}

func (q AppendQuery) Execute(database Database) (string, error) {
	value, exists := database.Data[q.Key]
	if exists {
		parsed, ok := value.([]string)
		if !ok {
			return "", errors.New("wrong type")
		} else {
			database.Data[q.Key] = append(parsed, q.Value)
		}
	} else {
		database.Data[q.Key] = []string{q.Value}
	}
	return "", nil
}

func ParseQuery(query string) Query {
	i := strings.Index(query, " ")
	if i < 0 {
		return EmptyQuery{}
	}
	operation := query[:i]
	key := query[i+1:]
	if operation == "get" {
		return GetQuery{key}
	}
	i = strings.Index(key, " ")
	if i < 0 {
		return EmptyQuery{}
	}
	value := key[i+1:]
	key = key[:i]
	if operation == "set" {
		return SetQuery{key, value}
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
