package main

import (
	"gotest.tools/v3/assert"
	"strings"
	"testing"
)

func TestExec(t *testing.T) {
	out, err := Exec("ls")
	assert.NilError(t, err)
	assert.Check(t, strings.Contains(out, "database.go"))
}
