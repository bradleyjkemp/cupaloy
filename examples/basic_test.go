package examples

import (
	"github.com/bradleyjkemp/cupaloy"
	"testing"
)

func TestString(t *testing.T) {
	result := "Hello world"
	err := cupaloy.Snapshot(result)
	if err != nil {
		t.Error("This will pass because \"Hello world\" is in the snapshot")
	}

	err = cupaloy.Snapshot("Hello world!")
	if err == nil {
		t.Error("Now it will fail because the snapshot doesn't have an exclamation mark")
	}
}

// Tests are independent of each other
func TestSecondString(t *testing.T) {
	result := "Hello Universe!"
	err := cupaloy.Snapshot(result)
	if err != nil {
		t.Error("This will pass because Snapshots are per test function")
	}
}

// Multiple snapshots can be taken in a single test
func TestMultipleStrings(t *testing.T) {
	result1 := "Hello"
	err := cupaloy.Snapshot(result1)
	if err != nil {
		t.Error("This will pass as normal")
	}

	result2 := "World"
	err = cupaloy.SnapshotMulti("result2", result2)
	if err != nil {
		t.Error("This will pass also as we've specified a unique (to this function) id")
	}
}
