package cupaloy

import (
	"fmt"
)

type Snapshotter interface {
	// Snapshot compares the given value to the it's previous value stored on the filesystem.
	// An error containing a diff is returned if the snapshots do not match.
	// Snapshot determines the snapshot file automatically from the name of the calling function.
	Snapshot(i ...interface{}) error

	// SnapshotMulti is identical to Snapshot but can be called multiple times from the same function.
	// This is done by providing a unique snapshotId for each invocation.
	SnapshotMulti(snapshotId string, i ...interface{}) error
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
func SnapshotMulti(snapshotId string, i ...interface{}) error {
	snapshotName := fmt.Sprintf("%s-%s", getNameOfCaller(), snapshotId)
	return defaultConfig().snapshot(snapshotName, i...)
}

func (c *config) Snapshot(i ...interface{}) error {
	return c.snapshot(getNameOfCaller(), i...)
}

func (c *config) SnapshotMulti(snapshotId string, i ...interface{}) error {
	snapshotName := fmt.Sprintf("%s-%s", getNameOfCaller(), snapshotId)
	return c.snapshot(snapshotName, i...)
}

func (c *config) snapshot(snapshotName string, i ...interface{}) error {
	snapshot := takeSnapshot(i...)

	if c.shouldUpdate() {
		err := c.writeSnapshot(snapshotName, snapshot)
		if err != nil {
			return err
		}

		return fmt.Errorf("snapshot updated for test %s", snapshotName)
	}

	prevSnapshot, err := c.readSnapshot(snapshotName)
	if err != nil {
		return err
	}

	if snapshot != prevSnapshot {
		diff := diffSnapshots(prevSnapshot, snapshot)
		return fmt.Errorf("snapshot not equal:\n%s\n", diff)
	}

	return nil
}
