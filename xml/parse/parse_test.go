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

	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
)

func TestElementsAndAttributesAndChars(t *testing.T) {
	astr := `<a k1="v1" k2="v2">b</a>`
	// [{"a": [{"k1": "v1"}, {"k2": v2"}, "b"]}]
	x := NewParser([]byte(astr))
	expect.Hint(t, x, parse.ArrayOpenHint)
	expect.Hint(t, x, parse.ObjectOpenHint)

	expect.Hint(t, x, parse.KeyHint)
	expect.String(t, x, "a")
	expect.Hint(t, x, parse.ArrayOpenHint)

	expect.Hint(t, x, parse.ObjectOpenHint)
	expect.Hint(t, x, parse.KeyHint)
	expect.String(t, x, "k1")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "v1")
	expect.Hint(t, x, parse.ObjectCloseHint)

	expect.Hint(t, x, parse.ObjectOpenHint)
	expect.Hint(t, x, parse.KeyHint)
	expect.String(t, x, "k2")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "v2")
	expect.Hint(t, x, parse.ObjectCloseHint)

	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "b")

	expect.Hint(t, x, parse.ArrayCloseHint)

	expect.Hint(t, x, parse.ObjectCloseHint)

	expect.Hint(t, x, parse.ArrayCloseHint)
	expect.EOF(t, x)
}

func TestParseElementsAndCharsAndWhiteSpaceManual(t *testing.T) {
	elemStr := `<A>
	<B>C</B>
</A>`
	// [{"A": ["\n\t", {"B": ["C"]}, "\n"]}]
	x := NewParser([]byte(elemStr))
	expect.Hint(t, x, parse.ArrayOpenHint)
	expect.Hint(t, x, parse.ObjectOpenHint)

	expect.Hint(t, x, parse.KeyHint)
	expect.String(t, x, "A")
	expect.Hint(t, x, parse.ArrayOpenHint)

	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "\n\t")

	expect.Hint(t, x, parse.ObjectOpenHint)
	expect.Hint(t, x, parse.KeyHint)
	expect.String(t, x, "B")

	expect.Hint(t, x, parse.ArrayOpenHint)
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "C")
	expect.Hint(t, x, parse.ArrayCloseHint)

	expect.Hint(t, x, parse.ObjectCloseHint)

	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "\n")

	expect.Hint(t, x, parse.ArrayCloseHint)

	expect.Hint(t, x, parse.ObjectCloseHint)
	expect.Hint(t, x, parse.ArrayCloseHint)
	expect.EOF(t, x)
}

func TestParseEmpty(t *testing.T) {
	str := ""
	// `[]`
	p := NewParser([]byte(str))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}

func TestParseChar(t *testing.T) {
	str := "a"
	// `[a]`
	p := NewParser([]byte(str))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}

func TestParseString(t *testing.T) {
	str := "abc"
	// `"abc"`
	p := NewParser([]byte(str))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "abc")
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}

func TestParseArray(t *testing.T) {
	str := "<a/><b/>"
	// `[a,b]`
	p := NewParser([]byte(str))
	expect.Hint(t, p, parse.ArrayOpenHint)

	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)

	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "b")
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)

	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}

func TestParseArrayNestedOpen(t *testing.T) {
	str := "<a><b/><c/></a>"
	// `[{"a":[{"b":[]},{"c":[]}]}]`
	p := NewParser([]byte(str))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ArrayOpenHint)

	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "b")
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)

	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "c")
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)

	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}
