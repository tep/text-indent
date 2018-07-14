package indent

import (
	"bytes"
	"fmt"
	"testing"
)

func TestIndent(t *testing.T) {
	var (
		inner = fmt.Sprintf("three [%c\nthis\nthat\nthose%c\n]", ShiftOut, ShiftIn)
		text  = fmt.Sprintf("one {%c\ntwo\n%s%c\n}\nfour", ShiftOut, inner, ShiftIn)
		in    = bytes.NewBufferString(text)
		out   = new(bytes.Buffer)
	)

	df := DefaultFilter

	df.Reader = in

	n, err := df.WriteTo(out)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Filtered %d bytes", n)
	t.Log("----------------------------------")
	t.Logf("%s\n", out.String())
	t.Log("----------------------------------")
}
