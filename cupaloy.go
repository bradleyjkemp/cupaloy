package cupaloy

import (
	"fmt"
)

// Config can be used to run cupaloy with customised behaviour e.g. changing how it decides to update snapshots
type Config struct {
	// ShouldUpdate allows you to control the decision of whether to update snapshots
	// The default behaviour is to check if the UPDATE_SNASPHOTS environment variable is set
	ShouldUpdate      func() bool
	subDirName        string
	snapshotExtension string
}

func DefaultConfig() *Config {
	return &Config{
		ShouldUpdate:      shouldUpdate,
		subDirName:        ".snapshots",
		snapshotExtension: "",
	}
}

// Snapshot compares the given value to the it's previous value stored on the filesystem.
// An error containing a diff is returned if the snapshots do not match.
// Snapshot determines the snapshot file automatically from the name of the calling function.
func Snapshot(i ...interface{}) error {
	return DefaultConfig().snapshot(getNameOfCaller(), i...)
}

// SnapshotMulti is identical to Snapshot but can be called multiple times from the same function.
// This is done by providing a unique snapshotId for each invocation.
func SnapshotMulti(snapshotId string, i ...interface{}) error {
	snapshotName := fmt.Sprintf("%s-%s", getNameOfCaller(), snapshotId)
	return DefaultConfig().snapshot(snapshotName, i...)
}

func (c *Config) Snapshot(i ...interface{}) error {
	return c.snapshot(getNameOfCaller(), i...)
}

func (c *Config) SnapshotMulti(snapshotId string, i ...interface{}) error {
	snapshotName := fmt.Sprintf("%s-%s", getNameOfCaller(), snapshotId)
	return c.snapshot(snapshotName, i...)
}

func (c *Config) snapshot(snapshotName string, i ...interface{}) error {
	snapshot := takeSnapshot(i...)

	if c.ShouldUpdate() {
		return c.writeSnapshot(snapshotName, snapshot)
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
