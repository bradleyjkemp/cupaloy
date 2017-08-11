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

// If a snapshot is update then this returns an error
// This is to prevent you accidentally updating your snapshots in CI
func TestUpdate(t *testing.T) {
	snapshotter := cupaloy.New(cupaloy.EnvVariableName("HOME"))

	err := snapshotter.Snapshot("Hello world")
	if err == nil {
		t.Errorf("This will always return an error %s", err)
	}
}

// If a snapshot doesn't exist then an error is thrown
func TestMissingSnapshot(t *testing.T) {
	snapshotter := cupaloy.New(cupaloy.EnvVariableName("ENOEXIST"))

	err := snapshotter.Snapshot("Hello world")
	if err == nil {
		t.Errorf("This will always return an error %s", err)
	}
}

// If the snapshots directory doesn't exit an error is returned
func TestMissingDirectory(t *testing.T) {
	snapshotter := cupaloy.New(
		cupaloy.EnvVariableName("ENOEXIST"),
		cupaloy.SnapshotSubdirectory("noexists"))

	err := snapshotter.Snapshot("Hello world")
	if err == nil {
		t.Errorf("This will always return an error %s", err)
	}
}
