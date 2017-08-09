package examples_test

import (
	"github.com/bradleyjkemp/cupaloy"
	"testing"
)

func TestString(t *testing.T) {
	result := "Hello advanced world!"
	err := cupaloy.Snapshot(result)
	if err != nil {
		t.Error("Tests in different packages are independent of each other")
	}
}

func TestConfig(t *testing.T) {
	snapshotter := cupaloy.DefaultConfig()
	snapshotter.ShouldUpdate = func() bool {
		return false
	}

	err := snapshotter.Snapshot("Hello Universe")
	if err != nil {
		t.Error("You can use a custom config struct to customise the behaviour of cupaloy")
	}

	err = snapshotter.SnapshotMulti("withExclamation", "Hello Universe!")
	if err != nil {
		t.Error("The config struct has all the same methods as the default")
	}
}

func TestMap(t *testing.T) {
	result := map[int]string{
		1: "Hello",
		3: "!",
		2: "World",
	}

	err := cupaloy.Snapshot(result)
	if err != nil {
		t.Error("Snapshots can be taken of any type")
	}
}
