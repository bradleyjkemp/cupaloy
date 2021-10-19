package cupaloy

import (
	"strings"
)

func replaceLineSeparator(snapshot string) string {
	return strings.ReplaceAll(snapshot, "\r\n", "\n")
}
