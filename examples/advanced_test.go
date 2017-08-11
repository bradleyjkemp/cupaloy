package examples_test

import (
	"github.com/bradleyjkemp/cupaloy"
	"testing"
)

// Snapshots are isolated by package so test functions with the same name are fine
func TestString(t *testing.T) {
	result := "Hello advanced world!"
	err := cupaloy.Snapshot(result)
	if err != nil {
		t.Error("Tests in different packages are independent of each other")
	}
}

// A configured instance of cupaloy has the same interface as the static methods
func TestConfig(t *testing.T) {
	snapshotter := cupaloy.New(cupaloy.EnvVariableName("UPDATE"))

	err := snapshotter.Snapshot("Hello Universe")
	if err != nil {
		t.Errorf("You can use a custom config struct to customise the behaviour of cupaloy %s", err)
	}

	err = snapshotter.SnapshotMulti("withExclamation", "Hello", "Universe!")
	if err != nil {
		t.Errorf("The config struct has all the same methods as the default %s", err)
	}
}

// All types can be snapshotted. Maps are snapshotted in a deterministic way
func TestMap(t *testing.T) {
	result := map[int]string{
		1: "Hello",
		3: "!",
		2: "World",
	}

	err := cupaloy.Snapshot(result)
	if err != nil {
		t.Errorf("Snapshots can be taken of any type %s", err)
	}
}
