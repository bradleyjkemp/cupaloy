package cupaloy

import (
	"errors"
	"github.com/pmezard/go-difflib/difflib"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	subDirName        = "snapshots"
	snapshotExtension = ""
)

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

func snapshotFilePath(testName string) string {
	return filepath.Join(subDirName, testName+snapshotExtension)
}

func takeSnapshot(i ...interface{}) string {
	return spewConfig.Sdump(i)
}

func readSnapshot(testName string) (string, error) {
	snapshotFile := snapshotFilePath(testName)
	buf, err := ioutil.ReadFile(snapshotFile)

	if os.IsNotExist(err) {
		return "", errors.New("no snapshot exists for test " + testName)
	}

	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func writeSnapshot(testName string, snapshot string) error {
	// check that ./snapshots/ directory exists before writing snapshot
	err := os.MkdirAll(subDirName, os.ModePerm)
	if err != nil {
		return errors.New("could not create snapshots directory")
	}

	snapshotFile := snapshotFilePath(testName)
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
