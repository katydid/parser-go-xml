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

// These tests were adapted from the json skip tests

// If the kind '{' was returned by Next, then the whole object is skipped.
func TestSkipObjectOpen(t *testing.T) {
	// str := `{"a":1,"b":2}`
	str := `<a>1</a><b>2</b>`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.EnterHint)
	expect.NoErr(t, p.Skip)
	// skipped over "a":1,"b":2}
	expect.EOF(t, p)
}

func TestSkipObjectNestedOpen(t *testing.T) {
	// str := `{"a":{"b":1,"c":2}}`
	str := `<a><b>1</b><c>2</c></a>`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.EnterHint)
	expect.NoErr(t, p.Skip)
	// skipped over "b":1,"c":2}
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

// If an object value was just parsed, then the rest of the object is skipped.
func TestSkipObjectKey(t *testing.T) {
	// str := `{"a":1,"b":2,"c":3}`
	str := `<a>1</a><b>2</b><c>3</c>`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.NoErr(t, p.Skip)
	// skipped over "b":2,"c":3}
	expect.EOF(t, p)
}

func TestSkipObjectNestedKey(t *testing.T) {
	// str := `{"a":{"b":1,"c":2,"d":3}}`
	str := `<a><b>1</b><c>2</c><d>3</d></a>`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "b")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 1)
	expect.NoErr(t, p.Skip)
	// skipped over "c":2,"d":3}
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

// If a object key was just parsed, then that key's value is skipped.
func TestSkipObjectValue(t *testing.T) {
	// str := `{"a":1,"b":2}`
	str := `<a>1</a><b>2</b>`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "b")
	expect.NoErr(t, p.Skip)
	// skipped over 2
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipObjectRecursiveValue(t *testing.T) {
	// str := `{"a":1,"b":{"c":{"d":{"e":"f"},"g":[1,2]}}}`
	str := `<a>1</a><b><c><d><e>f</e></d><g>h</g></c></b>`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "b")
	expect.NoErr(t, p.Skip)
	// skipped over {"c":{"d":{"e":"f"},"g":[1,2]}}
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestSkipObjectDeepRecursiveValue(t *testing.T) {
	// str := `{"a":1,"b":{"c":{"d":{"e":"f"},"g":[1,2]}}}`
	str := `<a>1</a><b><c><d><e>f</e></d><g>h</g></c></b>`
	p := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "b")

	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "c")
	expect.NoErr(t, p.Skip)
	// skipped over {"d":{"e":"f"},"g":[1,2]}
	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}
