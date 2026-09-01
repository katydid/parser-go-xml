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

	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
)

func TestParseElementManual(t *testing.T) {
	elemStr := `<A>B</A>`
	// {"A": "B"}
	x := NewParser(WithBuffer([]byte(elemStr)))
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")

	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "B")

	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestParseElementsManual(t *testing.T) {
	elemStr := `<A><B>C</B></A>`
	// {"A": {"B": "C"}}
	x := NewParser(WithBuffer([]byte(elemStr)))
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "B")

	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "C")

	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestParseElementsAndCharsAndWhiteSpaceManual(t *testing.T) {
	elemStr := `<A>
	<B>C</B>
</A>`
	// {"A": {"B": "C"}}
	x := NewParser(WithBuffer([]byte(elemStr)), WithSkipSpace())
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "B")

	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "C")

	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestAttr(t *testing.T) {
	astr := `<A k1="v1"/>`
	x := NewParser(WithBuffer([]byte(astr)), WithSkipSpace())
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "k1")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "v1")

	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestElementsAndAttributesAndChars(t *testing.T) {
	astr := `<a k1="v1" k2="v2">b</a>`
	// {"a": {"k1": "v1", "k2": "v2", "b": {}}}
	x := NewParser(WithBuffer([]byte(astr)), WithSkipSpace())
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "a")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "k1")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "v1")

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "k2")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "v2")

	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "b")

	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestParseEmpty(t *testing.T) {
	str := ""
	// `{}`
	p := NewParser(WithBuffer([]byte(str)), WithSkipSpace())
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestParseArray(t *testing.T) {
	str := "<a/><b/>"
	// `{"a": {}, "b": {}}`
	p := NewParser(WithBuffer([]byte(str)), WithSkipSpace())
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.LeaveHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "b")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.LeaveHint)

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestParseArrayNestedOpen(t *testing.T) {
	str := "<a><b/><c/></a>"
	// `{"a": {"b": {}, "c": {}}}`
	p := NewParser(WithBuffer([]byte(str)), WithSkipSpace())
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")

	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "b")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.LeaveHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "c")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.LeaveHint)

	expect.Hint(t, p, parse.LeaveHint)

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}
