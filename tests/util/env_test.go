package util

import "testing"

import (
	"../../src/util"
)

func TestGetEnv(t *testing.T) {
	d := "default value"
	env := util.GetEnv("test", d)
	if env != "default value" {
		t.Errorf("TestGetEnv was incorrect, got: %s, want: %s.", env, d)
	}
}
