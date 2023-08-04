package main

type Database struct {
	Data map[string]any
}

func NewDatabase() Database {
	return Database{make(map[string]any)}
}
