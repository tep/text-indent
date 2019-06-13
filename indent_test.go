// Copyright © 2018 Timothy E. Peoples <eng@toolman.org>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

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
