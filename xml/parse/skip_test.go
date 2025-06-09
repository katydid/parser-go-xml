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

	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
)

func TestSkipEmpty(t *testing.T) {
	str := ""
	// `[]`
	p := NewParser([]byte(str))
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}

func TestSkipChar(t *testing.T) {
	str := "a"
	// `[a]`
	p := NewParser([]byte(str))
	expect.NoErr(t, p.Skip)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}

func TestSkipUnknownString(t *testing.T) {
	str := "abc"
	// `"abc"`
	p := NewParser([]byte(str))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.NoErr(t, p.Skip)
	expect.EOF(t, p)
}

// // If the kind '[' was returned by Next, then the whole array is skipped.
func TestSkipArrayOpen(t *testing.T) {
	str := "<a/><b/>"
	// `[a,b]`
	p := NewParser([]byte(str))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.NoErr(t, p.Skip)
	// skipped over <a/>,<b/>]
	expect.EOF(t, p)
}

func TestSkipArrayNestedOpen(t *testing.T) {
	str := "<a><b/><c/></a>"
	// `[{"a":[{"b":[]},{"c":[]}]}]`
	p := NewParser([]byte(str))
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")
	expect.Hint(t, p, parse.ArrayOpenHint)

	expect.NoErr(t, p.Skip)
	// skipped over <b/>,<c/>]
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}

// // If an array element was parsed, then the rest of the array is skipped.
func TestSkipArrayElement(t *testing.T) {
	str := `<a><b/><c/><d/></a>`
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

	expect.NoErr(t, p.Skip)
	// skipped over <c/>,<d/>]
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.EOF(t, p)
}
