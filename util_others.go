//go:build !windows
// +build !windows

package cupaloy

func replaceLineSeparator(snapshot string) string {
	return snapshot
}
