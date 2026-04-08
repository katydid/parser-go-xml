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

package index

import (
	"testing"

	xmlparse "github.com/katydid/parser-go-xml/xml/parse"
	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
)

func TestElementsAndAttributesAndChars(t *testing.T) {
	astr := `<a k1="v1" k2="v2">b</a>`
	// [{"a": [{"k1": "v1"}, {"k2": v2"}, "b"]}]
	// [0: {"a": [0: {"k1": "v1"}, 1: {"k2": v2"}, 2: "b"]}]
	x := WithIndexedArrays(xmlparse.NewParser([]byte(astr)))
	expect.Hint(t, x, parse.ArrayOpenHint)
	expect.Hint(t, x, parse.KeyHint)
	expect.Int(t, x, 0)
	expect.Hint(t, x, parse.ObjectOpenHint)

	expect.Hint(t, x, parse.KeyHint)
	expect.String(t, x, "a")

	expect.Hint(t, x, parse.ArrayOpenHint)

	expect.Hint(t, x, parse.KeyHint)
	expect.Int(t, x, 0)
	expect.Hint(t, x, parse.ObjectOpenHint)
	expect.Hint(t, x, parse.KeyHint)
	expect.String(t, x, "k1")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "v1")
	expect.Hint(t, x, parse.ObjectCloseHint)

	expect.Hint(t, x, parse.KeyHint)
	expect.Int(t, x, 1)
	expect.Hint(t, x, parse.ObjectOpenHint)
	expect.Hint(t, x, parse.KeyHint)
	expect.String(t, x, "k2")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "v2")
	expect.Hint(t, x, parse.ObjectCloseHint)

	expect.Hint(t, x, parse.KeyHint)
	expect.Int(t, x, 2)
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "b")

	expect.Hint(t, x, parse.ArrayCloseHint)

	expect.Hint(t, x, parse.ObjectCloseHint)

	expect.Hint(t, x, parse.ArrayCloseHint)
	expect.EOF(t, x)
}
