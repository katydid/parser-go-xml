// Copyright 2026 Walter Schulze
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parse

import (
	"testing"

	"katydid.org.za/go/parser-go/expect"
	"katydid.org.za/go/parser-go/parse"
)

func TestSkipValue(t *testing.T) {
	elemStr := `<A><B>C</B></A>`
	// {"A": {"B": "C"}}
	x := NewParser(WithBuffer([]byte(elemStr)))
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "B")
	expect.NoErr(t, x.Skip)

	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestSkipSingleValue(t *testing.T) {
	elemStr := `<A>B</A>`
	// {"A": {"B": "C"}}
	x := NewParser(WithBuffer([]byte(elemStr)))
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")
	expect.NoErr(t, x.Skip)

	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestSkipElement(t *testing.T) {
	elemStr := `<A><B><C>D</C></B><E/></A>`
	// elemStr := `{"A": {"B": {"C": "D"}}, "E": {}}`
	x := NewParser(WithBuffer([]byte(elemStr)))
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "B")
	expect.NoErr(t, x.Skip)
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "E")
	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestSkipElementAfterEnter(t *testing.T) {
	elemStr := `<A><B><C>D</C></B><E/></A>`
	// elemStr := `{"A": {"B": {"C": "D"}}, "E": {}}`
	x := NewParser(WithBuffer([]byte(elemStr)))
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "B")
	expect.Hint(t, x, parse.EnterHint)
	expect.NoErr(t, x.Skip)

	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "E")
	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestSkipAttributeValue(t *testing.T) {
	str := `<A B="C" D="E"/>`
	x := NewParser(WithBuffer([]byte(str)))

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "B")
	expect.NoErr(t, x.Skip)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "D")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "E")

	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestSkipInElemStart(t *testing.T) {
	str := `<A><B>C</B></A>`
	x := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")

	expect.NoErr(t, x.Skip)

	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestSkipInElemMid(t *testing.T) {
	str := `<A><B>C</B><D>E</D></A>`
	x := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "B")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "C")

	expect.NoErr(t, x.Skip)

	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestSkipInElemAttr(t *testing.T) {
	str := `<A b="c" d="e"><B>C</B><D>E</D></A>`
	x := NewParser(WithBuffer([]byte(str)))

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "b")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "c")

	expect.NoErr(t, x.Skip)

	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}
