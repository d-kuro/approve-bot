package cmd

import (
	"bytes"
	"fmt"
	"testing"
)

func TestRunVersionCmd(t *testing.T) {
	buf := &bytes.Buffer{}
	o := NewOption(buf, nil)
	runVersionCmd(o)

	out := fmt.Sprintf("version: %s (rev: %s)\n", Version, Revision)
	if buf.String() != out {
		t.Errorf("got: %v, want: %v", buf.String(), out)
	}
}
