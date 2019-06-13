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
	"io"
	"os"
	"strings"
	"sync"
)

const (
	BufSize = 2048

	ShiftOut byte = 0x0E
	ShiftIn  byte = 0x0F
)

type Filter struct {
	Reader  io.Reader
	Indent  string
	Prefix  string
	BufSize int
	buf     []byte
	pos     int
	cnt     int
	tot     int64
	ready   bool
}

var DefaultFilter = &Filter{
	Reader:  os.Stdin,
	Indent:  "    ",
	BufSize: BufSize,
	buf:     make([]byte, BufSize),
}

type initializer struct {
	sync.Mutex
}

var initr initializer

func (i *initializer) prepare(f *Filter) {
	i.Lock()
	defer i.Unlock()

	if f.ready {
		return
	}
	defer func() { f.ready = true }()

	if f.Reader == nil {
		f.Reader = DefaultFilter.Reader
	}

	if f.Indent == "" {
		f.Indent = DefaultFilter.Indent
	}

	if f.BufSize == 0 {
		f.BufSize = DefaultFilter.BufSize
	}

	if f.buf == nil {
		f.buf = make([]byte, f.BufSize)
	}
}

func (f *Filter) WriteTo(w io.Writer) (int64, error) {
	initr.prepare(f)
	err := f.filter(w, 0)
	if err == io.EOF {
		err = nil
	}
	return f.tot, err
}

func (f *Filter) refresh() error {
	if f.pos < f.cnt {
		return nil
	}

	n, err := f.Reader.Read(f.buf)
	if err != nil {
		return err
	}

	f.cnt = n
	f.pos = 0

	return nil
}

func (f *Filter) filter(w io.Writer, lvl int) error {
	for {
		if err := f.refresh(); err != nil {
			return err
		}

		for f.pos < f.cnt {
			b := f.buf[f.pos]
			f.pos++

			switch b {
			case '\n':
				n, err := w.Write([]byte("\n" + f.Prefix + strings.Repeat(f.Indent, lvl)))
				f.tot += int64(n)
				if err != nil {
					return err
				}

			case ShiftOut:
				if err := f.filter(w, lvl+1); err != nil {
					return err
				}

			case ShiftIn:
				if lvl > 0 {
					return nil
				}

			default:
				n, err := w.Write([]byte{b})
				f.tot += int64(n)
				if err != nil {
					return err
				}
			}
		}
	}
}
