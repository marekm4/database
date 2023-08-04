package main

import "strconv"

type Database struct {
	Values   map[string]string
	Counters map[string]int
	Lists    map[string][]string
}

func NewDatabase() Database {
	return Database{make(map[string]string), make(map[string]int), make(map[string][]string)}
}

func (d Database) Select(key string) []string {
	if value, ok := d.Values[key]; ok {
		return []string{value}
	}
	if counter, ok := d.Counters[key]; ok {
		return []string{strconv.Itoa(counter)}
	}
	if list, ok := d.Lists[key]; ok {
		return list
	}
	return []string{""}
}

func (d Database) Update(key string, value string) {
	d.Values[key] = value
}

func (d Database) Increment(key string, value int) {
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
