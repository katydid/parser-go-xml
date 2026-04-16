//  Copyright 2025 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package parse

import (
	"testing"
)

func TestSkipEmpty(t *testing.T) {
	str := ""
	// `[]`
	p := NewParser(WithBuffer([]byte(str)))
	expectNoErr(t, p.Skip)
	expectHint(t, p, ArrayCloseHint)
	expectEOF(t, p)
}

func TestSkipChar(t *testing.T) {
	str := "a"
	// `[a]`
	p := NewParser(WithBuffer([]byte(str)))
	expectNoErr(t, p.Skip)
	expectHint(t, p, ValueHint)
	expectString(t, p, "a")
	expectHint(t, p, ArrayCloseHint)
	expectEOF(t, p)
}

func TestSkipUnknownString(t *testing.T) {
	str := "abc"
	// `"abc"`
	p := NewParser(WithBuffer([]byte(str)))
	expectHint(t, p, ArrayOpenHint)
	expectNoErr(t, p.Skip)
	expectEOF(t, p)
}

// If the kind '[' was returned by Next, then the whole array is skipped.
func TestSkipArrayOpen(t *testing.T) {
	str := "<a/><b/>"
	// `[a,b]`
	p := NewParser(WithBuffer([]byte(str)))
	expectHint(t, p, ArrayOpenHint)
	expectNoErr(t, p.Skip)
	// skipped over <a/>,<b/>]
	expectEOF(t, p)
}

func TestSkipArrayNestedOpen(t *testing.T) {
	str := "<a><b/><c/></a>"
	// `[{"a":[{"b":[]},{"c":[]}]}]`
	p := NewParser(WithBuffer([]byte(str)))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "a")
	expectHint(t, p, ArrayOpenHint)

	expectNoErr(t, p.Skip)
	// skipped over <b/>,<c/>]
	expectHint(t, p, ObjectCloseHint)
	expectHint(t, p, ArrayCloseHint)
	expectEOF(t, p)
}

// If an array element was parsed, then the rest of the array is skipped.
func TestSkipArrayElement(t *testing.T) {
	str := `<a><b/><c/><d/></a>`
	// `[{"a":[{"b":[]},{"c":[]},{"d":[]}]}]`
	p := NewParser(WithBuffer([]byte(str)))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "a")
	expectHint(t, p, ArrayOpenHint)

	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "b")
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ArrayCloseHint)
	expectHint(t, p, ObjectCloseHint)

	expectNoErr(t, p.Skip)
	// skipped over <c/>,<d/>]
	expectHint(t, p, ObjectCloseHint)
	expectHint(t, p, ArrayCloseHint)
	expectEOF(t, p)
}

// If the kind '{' was returned by Next, then the whole object is skipped.
func TestSkipObjectOpen(t *testing.T) {
	str := `<a><b/></a>`
	// `[{"a":[{"b":[]}]}]`
	p := NewParser(WithBuffer([]byte(str)))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ObjectOpenHint)
	expectNoErr(t, p.Skip)
	// skipped over "a":[{"b":[]}]}
	expectHint(t, p, ArrayCloseHint)
	expectEOF(t, p)
}

func TestSkipObjectNestedOpen(t *testing.T) {
	str := `<a><b><c/></b></a>`
	// `[{"a":[{"b":[{"c":[]}]}]}]`
	p := NewParser(WithBuffer([]byte(str)))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "a")
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ObjectOpenHint)
	expectNoErr(t, p.Skip)
	// skipped over "b":[{"c":[]}]}
	expectHint(t, p, ArrayCloseHint)
	expectHint(t, p, ObjectCloseHint)
	expectHint(t, p, ArrayCloseHint)
	expectEOF(t, p)
}

// If a object key was just parsed, then that key's value is skipped.
func TestSkipObjectValueValue(t *testing.T) {
	str := `<a><b>c</b>d</a>`
	// `[{"a":[{"b":["c"]}, "d"]}]`
	p := NewParser(WithBuffer([]byte(str)))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "a")
	expectHint(t, p, ArrayOpenHint)

	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "b")
	expectNoErr(t, p.Skip)
	// skipped over ["c"]
	expectHint(t, p, ObjectCloseHint)

	expectHint(t, p, ValueHint)
	expectString(t, p, "d")

	expectHint(t, p, ArrayCloseHint)
	expectHint(t, p, ObjectCloseHint)
	expectHint(t, p, ArrayCloseHint)

	expectEOF(t, p)
}

// If a object key was just parsed, then that key's value is skipped.
func TestSkipObjectValueObject(t *testing.T) {
	str := `<a><b><c/></b>d</a>`
	// `[{"a":[{"b":[{"c":[]}]}, "d"]}]`
	p := NewParser(WithBuffer([]byte(str)))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "a")
	expectHint(t, p, ArrayOpenHint)

	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "b")
	expectNoErr(t, p.Skip)
	// skipped over [{"c":[]}]
	expectHint(t, p, ObjectCloseHint)

	expectHint(t, p, ValueHint)
	expectString(t, p, "d")

	expectHint(t, p, ArrayCloseHint)
	expectHint(t, p, ObjectCloseHint)
	expectHint(t, p, ArrayCloseHint)

	expectEOF(t, p)
}

func TestSkipObjectRecursiveValue(t *testing.T) {
	str := `<a><b><c>f<u/>k</c></b>d</a>`
	// `[{"a":[{"b":[{"c":[...]}]}, "d"]}]`
	p := NewParser(WithBuffer([]byte(str)))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "a")
	expectHint(t, p, ArrayOpenHint)

	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "b")
	expectNoErr(t, p.Skip)
	// skipped over [{"c":[]}]
	expectHint(t, p, ObjectCloseHint)

	expectHint(t, p, ValueHint)
	expectString(t, p, "d")

	expectHint(t, p, ArrayCloseHint)
	expectHint(t, p, ObjectCloseHint)
	expectHint(t, p, ArrayCloseHint)

	expectEOF(t, p)
}
