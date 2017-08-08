package cupaloy

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

var spewConfig = spew.ConfigState{
	Indent:                  "  ",
	SortKeys:                true, // maps should be spewed in a deterministic order
	DisablePointerAddresses: true, // don't spew the addresses of pointers
	SpewKeys:                true, // if unable to sort map keys then spew keys to strings and sort those
}

// Snapshot compares the given value to the it's previous value stored on the filesystem.
// An error containing a diff is returned if the snapshots do not match.
// Snapshot determines the snapshot file automatically from the name of the calling function.
func Snapshot(i ...interface{}) error {
	return snapshot(getNameOfCaller(), i)
}

// SnapshotMulti is identical to Snapshot but can be called multiple times from the same function.
// This is done by providing a unique snapshotId for each invocation.
func SnapshotMulti(snapshotId string, i ...interface{}) error {
	return snapshot(fmt.Sprintf("%s-%s", getNameOfCaller(), snapshotId), i)
}

func snapshot(snapshotName string, i ...interface{}) error {
	snapshot := takeSnapshot(i)

	if shouldUpdate() {
		return writeSnapshot(snapshotName, snapshot)
	}

	prevSnapshot, err := readSnapshot(snapshotName)
	if err != nil {
		return err
	}

	if snapshot != prevSnapshot {
		diff := diffSnapshots(prevSnapshot, snapshot)
		return fmt.Errorf("snapshot not equal:\n%s\n", diff)
	}

	return nil
}
