package examples_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

// Snapshots are isolated by package so test functions with the same name are fine
func TestString(t *testing.T) {
	result := "Hello advanced world!"
	err := cupaloy.Snapshot(result)
	if err != nil {
		t.Fatalf("Tests in different packages are independent of each other %s", err)
	}
}

func TestRawString(t *testing.T) {
	result := "Hello advanced world!"
	result2 := "Goodbye"
	err := cupaloy.New(cupaloy.RawOutput(true)).Snapshot(result, result2)
	if err != nil {
		t.Fatalf("Strings can be output in raw form %s", err)
	}
}

func TestRawReader(t *testing.T) {
	result := &bytes.Buffer{}
	result.WriteString("Hello world")
	result2 := map[string]bool{
		"Hello": true,
		"World": true,
	}

	err := cupaloy.New(cupaloy.RawOutput(true)).Snapshot(result, result2)
	if err != nil {
		t.Fatalf("io.Readers can be output in raw form but other types remain the same %s", err)
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

// If a snapshot is update then this returns an error
// This is to prevent you accidentally updating your snapshots in CI
func TestUpdate(t *testing.T) {
	snapshotter := cupaloy.New(cupaloy.EnvVariableName("GOPATH"))

	err := snapshotter.Snapshot("Hello world")
	if err == nil {
		t.Fatalf("This will always return an error %s", err)
	}
	if err.Error() != "snapshot updated for test examples_test-TestUpdate" {
		t.Fatalf("Error returned will say that snapshot was updated")
	}
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
