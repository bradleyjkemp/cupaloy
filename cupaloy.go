package cupaloy

import (
	"fmt"
	"os"
	"strings"

	"github.com/bradleyjkemp/cupaloy/v2/internal"
)

// New constructs a new, configured instance of cupaloy using the given
// Configurators applied to the default config.
func New(configurators ...Configurator) *Config {
	return NewDefaultConfig().WithOptions(configurators...)
}

// Snapshot calls Snapshotter.Snapshot with the global config.
func Snapshot(i ...interface{}) error {
	return Global.snapshot(getNameOfCaller(), i...)
}

// SnapshotMulti calls Snapshotter.SnapshotMulti with the global config.
func SnapshotMulti(snapshotID string, i ...interface{}) error {
	snapshotName := fmt.Sprintf("%s-%s", getNameOfCaller(), snapshotID)
	return Global.snapshot(snapshotName, i...)
}

// SnapshotT calls Snapshotter.SnapshotT with the global config.
func SnapshotT(t TestingT, i ...interface{}) {
	t.Helper()
	Global.SnapshotT(t, i...)
}

// SnapshotWithName calls Snapshotter.SnapshotWithName with the global config.
func SnapshotWithName(snapshotName string, i ...interface{}) error {
	return Global.SnapshotWithName(snapshotName, i...)
}

// Snapshot compares the given variable to its previous value stored on the filesystem.
// An error containing a diff is returned if the snapshots do not match, or if a new
// snapshot was created.
//
// Snapshot determines the snapshot file automatically from the name of the calling function.
// As a result it can be called at most once per function. If you want to call Snapshot
// multiple times in a function, if possible, instead collect the values and call Snapshot
// with all values at once. Otherwise see SnapshotMulti.
//
// If using snapshots in tests, prefer the SnapshotT function which fails the test
// directly, rather than requiring your to remember to check the error.
func (c *Config) Snapshot(i ...interface{}) error {
	return c.snapshot(getNameOfCaller(), i...)
}

// SnapshotMulti is similar to Snapshot but can be called multiple times from the
// same function. This is possible by providing a unique id for each snapshot which is
// appended to the function name to form the snapshot name.
func (c *Config) SnapshotMulti(snapshotID string, i ...interface{}) error {
	snapshotName := fmt.Sprintf("%s-%s", getNameOfCaller(), snapshotID)
	return c.snapshot(snapshotName, i...)
}

// SnapshotWithName is similar to SnapshotMulti without appending the function name.
// It is useful when you need full control of the snapshot filename.
func (c *Config) SnapshotWithName(snapshotName string, i ...interface{}) error {
	return c.snapshot(snapshotName, i...)
}

// SnapshotT compares the given variable to the its previous value stored on the filesystem.
// The current test is failed (with error containing a diff) if the values do not match, or
// if a new snapshot was created.
//
// SnapshotT determines the snapshot file automatically from the name of the test (using
// the t.Name() function). As a result, SnapshotT can be called at most once per test.
// If you want to call SnapshotT multiple times in a test, if possible, instead collect the
// values and call SnapshotT with all values at once. Alternatively, use sub-tests and call
// SnapshotT once in each.
//
// If using snapshots in tests, SnapshotT is preferred over Snapshot and SnapshotMulti.
func (c *Config) SnapshotT(t TestingT, i ...interface{}) {
	t.Helper()
	if t.Failed() {
		return
	}

	snapshotName := strings.Replace(t.Name(), "/", "-", -1)
	err := c.snapshot(snapshotName, i...)
	if err != nil {
		if c.fatalOnMismatch {
			t.Fatal(err)
			return
		}
		t.Error(err)
	}
}

// WithOptions returns a copy of an existing Config with additional Configurators applied.
// This can be used to apply a different option for a single call e.g.
//  snapshotter.WithOptions(cupaloy.SnapshotSubdirectory("testdata")).SnapshotT(t, result)
// Or to modify the Global Config e.g.
//  cupaloy.Global = cupaloy.Global.WithOptions(cupaloy.SnapshotSubdirectory("testdata"))
func (c *Config) WithOptions(configurators ...Configurator) *Config {
	clonedConfig := c.clone()

	for _, configurator := range configurators {
		configurator(clonedConfig)
	}

	return clonedConfig
}

func (c *Config) snapshot(snapshotName string, i ...interface{}) error {
	snapshot := c.takeSnapshot(i...)

	prevSnapshot, err := c.readSnapshot(snapshotName)
	if os.IsNotExist(err) {
		if c.createNewAutomatically {
			return c.updateSnapshot(snapshotName, prevSnapshot, snapshot)
		}
		return internal.ErrNoSnapshot{Name: snapshotName}
	}
	if err != nil {
		return err
	}

	if snapshot == prevSnapshot || c.takeV1Snapshot(i...) == prevSnapshot {
		// previous snapshot matches current value
		return nil
	}

	if c.shouldUpdate() {
		// updates snapshot to current value and upgrades snapshot format
		return c.updateSnapshot(snapshotName, prevSnapshot, snapshot)
	}

	return internal.ErrSnapshotMismatch{
		Diff: diffSnapshots(prevSnapshot, snapshot),
	}
}
