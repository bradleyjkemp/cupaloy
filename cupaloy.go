package cupaloy

import (
	"fmt"
	"os"
)

// Snapshotter is the API for taking snapshots of values in your tests.
type Snapshotter interface {
	// Snapshot compares the given value to the it's previous value stored on the filesystem.
	// An error containing a diff is returned if the snapshots do not match.
	// Snapshot determines the snapshot file automatically from the name of the calling function.
	Snapshot(i ...interface{}) error

	// SnapshotMulti is identical to Snapshot but can be called multiple times from the same function.
	// This is done by providing a unique snapshotId for each invocation.
	SnapshotMulti(snapshotID string, i ...interface{}) error
}

// New constructs a new, configured instance of cupaloy using the given Configurators.
func New(configurators ...Configurator) Snapshotter {
	config := defaultConfig()

	for _, configurator := range configurators {
		configurator(config)
	}

	return config
}

// Snapshot calls Snapshotter.Snapshot with the default config.
func Snapshot(i ...interface{}) error {
	return defaultConfig().snapshot(getNameOfCaller(), i...)
}

// SnapshotMulti calls Snapshotter.SnapshotMulti with the default config.
func SnapshotMulti(snapshotID string, i ...interface{}) error {
	snapshotName := fmt.Sprintf("%s-%s", getNameOfCaller(), snapshotID)
	return defaultConfig().snapshot(snapshotName, i...)
}

func (c *config) Snapshot(i ...interface{}) error {
	return c.snapshot(getNameOfCaller(), i...)
}

func (c *config) SnapshotMulti(snapshotID string, i ...interface{}) error {
	snapshotName := fmt.Sprintf("%s-%s", getNameOfCaller(), snapshotID)
	return c.snapshot(snapshotName, i...)
}

func (c *config) snapshot(snapshotName string, i ...interface{}) error {
	snapshot := takeSnapshot(i...)

	prevSnapshot, err := c.readSnapshot(snapshotName)
	if os.IsNotExist(err) {
		return c.updateSnapshot(snapshotName, snapshot)
	}
	if err != nil {
		return err
	}

	if c.shouldUpdate() {
		return c.updateSnapshot(snapshotName, snapshot)
	}

	if snapshot != prevSnapshot {
		diff := diffSnapshots(prevSnapshot, snapshot)
		return fmt.Errorf("snapshot not equal:\n%s", diff)
	}

	return nil
}
