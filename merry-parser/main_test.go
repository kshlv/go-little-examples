package main

import (
	"reflect"
	"testing"
)

func TestIndex(t *testing.T) {
	type josh struct {
		Name string
		Age  int
	}
	j := josh{"Josh", 15}

	_, actualIndex := Index(j)
	expectedIndex := "0,24,1000;2,6,1000001;9,13,1000010;16,19,1000001;21,23,100010"

	if actualIndex != expectedIndex {
		t.Errorf("expected index: %q, actual index: %q", expectedIndex, actualIndex)
	}
}

func TestJSON(t *testing.T) {
	type josh struct {
		Name string
		Age  int
	}
	j := josh{"Josh", 15}

	actualJSON, _ := Index(j)
	expectedJSON := []byte("{\"name\":\"Josh\",\"age\":15}")

	if !reflect.DeepEqual(actualJSON, expectedJSON) {
		t.Errorf("expected JSON: %s, actual JSON: %s", expectedJSON, actualJSON)
	}
}
