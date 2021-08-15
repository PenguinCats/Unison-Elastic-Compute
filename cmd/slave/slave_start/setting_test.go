package main

import "testing"

func TestLoadGlobalSetting(t *testing.T) {
	err := LoadGlobalSetting("./slave.ini")
	if err != nil {
		t.Fatalf("load config fail")
	}
}
