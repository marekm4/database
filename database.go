package main

import (
	"errors"
	"strconv"
	"strings"
)

type Database struct {
	Data map[string]any
}

func NewDatabase() Database {
	return Database{make(map[string]any)}
}

type Query interface {
	Execute(Database) (string, error)
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
		numericValue, err := strconv.Atoi(value)
		if err != nil {
			return nil, errors.New("invalid numeric value")
		}
		return IncrementQuery{key, numericValue}, nil
	}
	if operation == "append" {
		return AppendQuery{key, value}, nil
	}
	return nil, errors.New("invalid operation")
}
