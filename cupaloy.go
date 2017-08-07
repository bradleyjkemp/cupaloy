package cupaloy

import (
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

var spewConfig = spew.ConfigState{
	Indent:                  "  ",
	SortKeys:                true, // maps should be spewed in a deterministic order
	DisablePointerAddresses: true, // don't spew the addresses of pointers
	SpewKeys:                true, // if unable to sort map keys then spew keys to strings and sort those
}

func Snapshot(i ...interface{}) error {
	return snapshot(getNameOfCaller(), i)
}

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
		return errors.New(fmt.Sprintf("snapshot not equal:\n%s\n", diff))
	}

	return nil
}
