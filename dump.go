package main

import (
	"bufio"
	"fmt"
	"os"
)

func Dump(database Database, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	for _, query := range DumpQueries(database) {
		_, err = file.WriteString(query + "\n")
		if err != nil {
			return err
		}
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func DumpQueries(database Database) []string {
	queries := []string{}
	for key, value := range database.Values {
		queries = append(queries, fmt.Sprintf("update %s %s", key, value))
	}
	for key, value := range database.Counters {
		queries = append(queries, fmt.Sprintf("increment %s %d", key, value))
	}
	for key, values := range database.Lists {
		for _, value := range values {
			queries = append(queries, fmt.Sprintf("append %s %s", key, value))
		}
	}
	return queries
}

func Load(database Database, filename string) error {
	if _, err := os.Stat(filename); err != nil {
		return nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	queries := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		queries = append(queries, line)
	}
	database.Clear()
	LoadQueries(database, queries)

	err = scanner.Err()
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func LoadQueries(database Database, queries []string) {
	for _, queryString := range queries {
		query := ParseQuery(queryString)
		query.Execute(database)
	}
}
