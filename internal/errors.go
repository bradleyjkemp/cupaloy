package internal

import "fmt"

type ErrSnapshotCreated struct {
	Name     string
	Contents string
}

func (e ErrSnapshotCreated) Error() string {
	return fmt.Sprintf("snapshot created for test %s, with contents:\n%s", e.Name, e.Contents)
}

type ErrSnapshotUpdated struct {
	Name string
	Diff string
}

func (e ErrSnapshotUpdated) Error() string {
	return fmt.Sprintf("snapshot %s updated:\n%s", e.Name, e.Diff)
}

type ErrSnapshotMismatch struct {
	Diff string
}

func (e ErrSnapshotMismatch) Error() string {
	return fmt.Sprintf("snapshot not equal:\n%s", e.Diff)
}

type ErrNoSnapshot struct {
	Name string
}

func (e ErrNoSnapshot) Error() string {
	return fmt.Sprintf("snapshot %s does not exist", e.Name)
}
