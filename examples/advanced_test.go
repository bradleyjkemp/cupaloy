package examples_test

import (
	"github.com/bradleyjkemp/cupaloy"
	"testing"
)

func TestString(t *testing.T) {
	result := "Hello advanced world!"
	err := cupaloy.Snapshot(result)
	if err != nil {
		t.Error("Tests in different packages are independent of each other")
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
