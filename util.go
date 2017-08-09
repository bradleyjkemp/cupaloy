package cupaloy

import (
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/pmezard/go-difflib/difflib"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var spewConfig = spew.ConfigState{
	Indent:                  "  ",
	SortKeys:                true, // maps should be spewed in a deterministic order
	DisablePointerAddresses: true, // don't spew the addresses of pointers
	DisableCapacities:       true, // don't spew capacities of collections
	SpewKeys:                true, // if unable to sort map keys then spew keys to strings and sort those
}

func getNameOfCaller() string {
	pc, _, _, _ := runtime.Caller(2) // first caller is the caller of this function, we want the caller of our caller
	fullPath := runtime.FuncForPC(pc).Name()
	packageFunctionName := filepath.Base(fullPath)

	return strings.Replace(packageFunctionName, ".", "-", -1)
}

func shouldUpdate() bool {
	_, varSet := os.LookupEnv("UPDATE_SNAPSHOTS")
	return varSet
}

func (c *Config) snapshotFilePath(testName string) string {
	return filepath.Join(c.subDirName, testName+c.snapshotExtension)
}

func takeSnapshot(i ...interface{}) string {
	return spewConfig.Sdump(i...)
}

func (c *Config) readSnapshot(snapshotName string) (string, error) {
	snapshotFile := c.snapshotFilePath(snapshotName)
	buf, err := ioutil.ReadFile(snapshotFile)

	if os.IsNotExist(err) {
		return "", fmt.Errorf("no snapshot exists for test %s", snapshotName)
	}

	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (c *Config) writeSnapshot(snapshotName string, snapshot string) error {
	// check that subdirectory exists before writing snapshot
	err := os.MkdirAll(c.subDirName, os.ModePerm)
	if err != nil {
		return errors.New("could not create snapshots directory")
	}

	snapshotFile := c.snapshotFilePath(snapshotName)
	return ioutil.WriteFile(snapshotFile, []byte(snapshot), os.FileMode(0644))
}

func diffSnapshots(previous, current string) string {
	diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(previous),
		B:        difflib.SplitLines(current),
		FromFile: "Previous",
		FromDate: "",
		ToFile:   "Current",
		ToDate:   "",
		Context:  1,
	})

	return diff
}
