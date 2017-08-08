package advanced

import (
	"github.com/bradleyjkemp/cupaloy"
	"testing"
)

func TestSubfolder(t *testing.T) {
	result := "Hello world"
	err := cupaloy.Snapshot(result)
	if err != nil {
		t.Error("Snapshots are stored relative to the test they are taken in")
	}
}

func TestMap(t *testing.T) {
	result := map[int]string{
		1: "Hello",
		3: "!",
		2: "World",
	}

	err := cupaloy.Snapshot(result)
	if err != nil {
		t.Error("Snapshots can be taken of any type")
	}
}
