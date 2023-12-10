package main

import (
	"fmt"
	"strings"
)

type Database struct {
	Values   map[string]string
	Counters map[string]float64
	Lists    map[string][]string
}

func NewDatabase() Database {
	return Database{make(map[string]string), make(map[string]float64), make(map[string][]string)}
}

func (d Database) List(prefix string) []string {
	keys := []string{}
	for key := range d.Values {
		if strings.HasPrefix(key, prefix) {
			keys = append(keys, key)
		}
	}
	for key := range d.Counters {
		if strings.HasPrefix(key, prefix) {
			keys = append(keys, key)
		}
	}
	for key := range d.Lists {
		if strings.HasPrefix(key, prefix) {
			keys = append(keys, key)
		}
	}
	return keys
}

func (d Database) Select(key string) []string {
	if value, ok := d.Values[key]; ok {
		return []string{value}
	}
	if counter, ok := d.Counters[key]; ok {
		return []string{fmt.Sprintf("%f", counter)}
	}
	if list, ok := d.Lists[key]; ok {
		return list
	}
	return []string{""}
}

func (d Database) Update(key string, value string) {
	d.Values[key] = value
}

func (d Database) Increment(key string, value float64) {
	if counter, ok := d.Counters[key]; ok {
		d.Counters[key] = counter + value
	} else {
		d.Counters[key] = value
	}
}

func (d Database) Append(key string, value string) {
	d.Lists[key] = append(d.Lists[key], value)
}

func (d Database) Clear() {
	for key := range d.Values {
		delete(d.Values, key)
	}
	for key := range d.Counters {
		delete(d.Counters, key)
	}
	for key := range d.Lists {
		delete(d.Lists, key)
	}
}
