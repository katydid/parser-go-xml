// Copyright 2025 Walter Schulze
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
)

func TestElementsAndAttributesAndChars(t *testing.T) {
	astr := `<a k1="v1" k2="v2">b</a>`
	// [{"a": [{"k1": "v1"}, {"k2": v2"}, "b"]}]
	x := NewParser([]byte(astr))
	expectHint(t, x, ArrayOpenHint)
	expectHint(t, x, ObjectOpenHint)

	expectHint(t, x, KeyHint)
	expectString(t, x, "a")
	expectHint(t, x, ArrayOpenHint)

	expectHint(t, x, ObjectOpenHint)
	expectHint(t, x, KeyHint)
	expectString(t, x, "k1")
	expectHint(t, x, ValueHint)
	expectString(t, x, "v1")
	expectHint(t, x, ObjectCloseHint)

	expectHint(t, x, ObjectOpenHint)
	expectHint(t, x, KeyHint)
	expectString(t, x, "k2")
	expectHint(t, x, ValueHint)
	expectString(t, x, "v2")
	expectHint(t, x, ObjectCloseHint)

	expectHint(t, x, ValueHint)
	expectString(t, x, "b")

	expectHint(t, x, ArrayCloseHint)

	expectHint(t, x, ObjectCloseHint)

	expectHint(t, x, ArrayCloseHint)
	expectEOF(t, x)
}

func TestParseElementsAndCharsAndWhiteSpaceManual(t *testing.T) {
	elemStr := `<A>
	<B>C</B>
</A>`
	// [{"A": ["\n\t", {"B": ["C"]}, "\n"]}]
	x := NewParser([]byte(elemStr))
	expectHint(t, x, ArrayOpenHint)
	expectHint(t, x, ObjectOpenHint)

	expectHint(t, x, KeyHint)
	expectString(t, x, "A")
	expectHint(t, x, ArrayOpenHint)

	expectHint(t, x, ValueHint)
	expectString(t, x, "\n\t")

	expectHint(t, x, ObjectOpenHint)
	expectHint(t, x, KeyHint)
	expectString(t, x, "B")

	expectHint(t, x, ArrayOpenHint)
	expectHint(t, x, ValueHint)
	expectString(t, x, "C")
	expectHint(t, x, ArrayCloseHint)

	expectHint(t, x, ObjectCloseHint)

	expectHint(t, x, ValueHint)
	expectString(t, x, "\n")

	expectHint(t, x, ArrayCloseHint)

	expectHint(t, x, ObjectCloseHint)
	expectHint(t, x, ArrayCloseHint)
	expectEOF(t, x)
}

func TestParseEmpty(t *testing.T) {
	str := ""
	// `[]`
	p := NewParser([]byte(str))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ArrayCloseHint)
	expectEOF(t, p)
}

func TestParseChar(t *testing.T) {
	str := "a"
	// `[a]`
	p := NewParser([]byte(str))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ValueHint)
	expectString(t, p, "a")
	expectHint(t, p, ArrayCloseHint)
	expectEOF(t, p)
}

func TestParseString(t *testing.T) {
	str := "abc"
	// `"abc"`
	p := NewParser([]byte(str))
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ValueHint)
	expectString(t, p, "abc")
	expectHint(t, p, ArrayCloseHint)
	expectEOF(t, p)
}

func TestParseArray(t *testing.T) {
	str := "<a/><b/>"
	// `[a,b]`
	p := NewParser([]byte(str))
	expectHint(t, p, ArrayOpenHint)

	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "a")
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ArrayCloseHint)
	expectHint(t, p, ObjectCloseHint)

	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "b")
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ArrayCloseHint)
	expectHint(t, p, ObjectCloseHint)

	expectHint(t, p, ArrayCloseHint)
	expectEOF(t, p)
}

func TestParseArrayNestedOpen(t *testing.T) {
	str := "<a><b/><c/></a>"
	// `[{"a":[{"b":[]},{"c":[]}]}]`
	p := NewParser([]byte(str))
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

	expectHint(t, p, ObjectOpenHint)
	expectHint(t, p, KeyHint)
	expectString(t, p, "c")
	expectHint(t, p, ArrayOpenHint)
	expectHint(t, p, ArrayCloseHint)
	expectHint(t, p, ObjectCloseHint)

	expectHint(t, p, ArrayCloseHint)
	expectHint(t, p, ObjectCloseHint)
	expectHint(t, p, ArrayCloseHint)
	expectEOF(t, p)
}
