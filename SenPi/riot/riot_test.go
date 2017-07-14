package riot

import (
	"strings"
	"testing"
)

func TestParseCoAPResourceResponse(t *testing.T) {
	var samplePaths []string
	samplePaths = append(samplePaths, "/path1")
	samplePaths = append(samplePaths, "/path2")
	samplePaths = append(samplePaths, "/path3")
	joinedPaths := "<" + strings.Join(samplePaths, ">,<") + ">"

	parsedPaths := ParseCoAPResourceResponse([]byte(joinedPaths))

	if !isArrayEqual(samplePaths, parsedPaths) {
		t.Fatalf("parsed paths are not equal")
	}
}

func isArrayEqual(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
