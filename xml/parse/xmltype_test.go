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

func expectXMLType(t *testing.T, p Parser, want XMLType) {
	t.Helper()
	got := p.TokenXMLType()
	if got != want {
		t.Fatalf("expected xml type %c, but got %c", want, got)
	}
}

func TestXMLType(t *testing.T) {
	astr := `<a k1="v1" k2="v2">b</a>`
	// [{"a": [{"k1": "v1"}, {"k2": v2"}, "b"]}]
	x := NewParser([]byte(astr))
	expect.Hint(t, x, parse.ArrayOpenHint)
	expect.Hint(t, x, parse.ObjectOpenHint)

	expect.Hint(t, x, parse.KeyHint)
	expectXMLType(t, x, ElemXMLType)
	expect.String(t, x, "a")
	expect.Hint(t, x, parse.ArrayOpenHint)

	expect.Hint(t, x, parse.ObjectOpenHint)
	expect.Hint(t, x, parse.KeyHint)
	expectXMLType(t, x, AttrXMLType)
	expect.String(t, x, "k1")
	expect.Hint(t, x, parse.ValueHint)
	expectXMLType(t, x, TextXMLType)
	expect.String(t, x, "v1")
	expect.Hint(t, x, parse.ObjectCloseHint)

	expect.Hint(t, x, parse.ObjectOpenHint)
	expect.Hint(t, x, parse.KeyHint)
	expectXMLType(t, x, AttrXMLType)
	expect.String(t, x, "k2")
	expect.Hint(t, x, parse.ValueHint)
	expectXMLType(t, x, TextXMLType)
	expect.String(t, x, "v2")
	expect.Hint(t, x, parse.ObjectCloseHint)

	expect.Hint(t, x, parse.ValueHint)
	expectXMLType(t, x, TextXMLType)
	expect.String(t, x, "b")

	expect.Hint(t, x, parse.ArrayCloseHint)

	expect.Hint(t, x, parse.ObjectCloseHint)

	expect.Hint(t, x, parse.ArrayCloseHint)
}
