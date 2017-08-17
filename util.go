package cupaloy

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pmezard/go-difflib/difflib"
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

func envVariableSet(envVariable string) bool {
	_, varSet := os.LookupEnv(envVariable)
	return varSet
}

func (c *config) snapshotFilePath(testName string) string {
	return filepath.Join(c.subDirName, testName+c.snapshotExtension)
}

func takeSnapshot(i ...interface{}) string {
	return spewConfig.Sdump(i...)
}

func (c *config) readSnapshot(snapshotName string) (string, error) {
	snapshotFile := c.snapshotFilePath(snapshotName)
	buf, err := ioutil.ReadFile(snapshotFile)

	if os.IsNotExist(err) {
		return "", err
	}

	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (c *config) updateSnapshot(snapshotName string, snapshot string) error {
	// check that subdirectory exists before writing snapshot
	err := os.MkdirAll(c.subDirName, os.ModePerm)
	if err != nil {
		return errors.New("could not create snapshots directory")
	}

	snapshotFile := c.snapshotFilePath(snapshotName)
	_, err = os.Stat(snapshotFile)
	isNewSnapshot := os.IsNotExist(err)

	err = ioutil.WriteFile(snapshotFile, []byte(snapshot), os.FileMode(0644))
	if err != nil {
		return err
	}

	if isNewSnapshot {
		return fmt.Errorf("snapshot created for test %s", snapshotName)
	} else {
		return fmt.Errorf("snapshot updated for test %s", snapshotName)
	}
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
