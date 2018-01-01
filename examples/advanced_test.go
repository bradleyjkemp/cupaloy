package examples_test

import (
	"io/ioutil"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

// Snapshots are isolated by package so test functions with the same name are fine
func TestString(t *testing.T) {
	result := "Hello advanced world!"
	err := cupaloy.Snapshot(result)
	if err != nil {
		t.Fatal("Tests in different packages are independent of each other")
	}
}

// A configured instance of cupaloy has the same interface as the static methods
func TestConfig(t *testing.T) {
	snapshotter := cupaloy.New(cupaloy.EnvVariableName("UPDATE"))

	err := snapshotter.Snapshot("Hello Universe")
	if err != nil {
		t.Fatalf("You can use a custom config struct to customise the behaviour of cupaloy %s", err)
	}

	err = snapshotter.SnapshotMulti("withExclamation", "Hello", "Universe!")
	if err != nil {
		t.Fatalf("The config struct has all the same methods as the default %s", err)
	}
}

// If a snapshot is updated then this returns an error
// This is to prevent you accidentally updating your snapshots in CI
func TestUpdate(t *testing.T) {
	snapshotter := cupaloy.New(cupaloy.EnvVariableName("GOPATH"))

	err := snapshotter.Snapshot("Hello world")
	if err != nil {
		t.Fatalf("Updating a snapshot with the same value does not fail a test %s", err)
	}

	err = snapshotter.Snapshot("Hello new world")
	if err == nil {
		t.Fatalf("Updating a snapshot with a new value is always an error %s", err)
	}
	if err.Error() != "snapshot updated for test examples_test-TestUpdate" {
		t.Fatalf("Error returned will say that snapshot was updated")
	}

	snapshotter.Snapshot("Hello world") // reset snapshot to known state
}

// If a snapshot doesn't exist then it is created and an error returned
func TestMissingSnapshot(t *testing.T) {
	tempdir, err := ioutil.TempDir(".", "ignored")
	if err != nil {
		t.Fatal(err)
	}

	snapshotter := cupaloy.New(
		cupaloy.EnvVariableName("ENOEXIST"),
		cupaloy.SnapshotSubdirectory(tempdir))

	err = snapshotter.Snapshot("Hello world")
	if err == nil {
		t.Fatalf("This will always return an error %s", err)
	}
	if err.Error() != "snapshot created for test examples_test-TestMissingSnapshot" {
		t.Fatalf("Error returned will say that snapshot was created %s", err)
	}
}

// Multiple snapshots can be taken in a single test
func TestMultipleSnapshots(t *testing.T) {
	t.Run("hello", func(t *testing.T) {
		result1 := "Hello"
		cupaloy.SnapshotT(t, result1)
	})

	t.Run("world", func(t *testing.T) {
		result2 := "World"
		cupaloy.New().SnapshotT(t, result2)
	})
}
