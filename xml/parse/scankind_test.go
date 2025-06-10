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

	"github.com/katydid/parser-go-xml/xml/scan"
	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
)

func expectScanKind(t *testing.T, p Parser, want scan.Kind) {
	t.Helper()
	got := p.ScanKind()
	if got != want {
		t.Fatalf("expected xml type %c, but got %c", want, got)
	}
}

func TestScanKind(t *testing.T) {
	astr := `<a k1="v1" k2="v2">b</a>`
	// [{"a": [{"k1": "v1"}, {"k2": v2"}, "b"]}]
	x := NewParser([]byte(astr))
	expect.Hint(t, x, parse.ArrayOpenHint)
	expect.Hint(t, x, parse.ObjectOpenHint)

	expect.Hint(t, x, parse.KeyHint)
	expectScanKind(t, x, scan.StartKind)
	expect.String(t, x, "a")
	expect.Hint(t, x, parse.ArrayOpenHint)

	expect.Hint(t, x, parse.ObjectOpenHint)
	expect.Hint(t, x, parse.KeyHint)
	expectScanKind(t, x, scan.AttrKeyKind)
	expect.String(t, x, "k1")
	expect.Hint(t, x, parse.ValueHint)
	expectScanKind(t, x, scan.AttrValueKind)
	expect.String(t, x, "v1")
	expect.Hint(t, x, parse.ObjectCloseHint)

	expect.Hint(t, x, parse.ObjectOpenHint)
	expect.Hint(t, x, parse.KeyHint)
	expectScanKind(t, x, scan.AttrKeyKind)
	expect.String(t, x, "k2")
	expect.Hint(t, x, parse.ValueHint)
	expectScanKind(t, x, scan.AttrValueKind)
	expect.String(t, x, "v2")
	expect.Hint(t, x, parse.ObjectCloseHint)

	expect.Hint(t, x, parse.ValueHint)
	expectScanKind(t, x, scan.CharKind)
	expect.String(t, x, "b")

	expect.Hint(t, x, parse.ArrayCloseHint)

	expect.Hint(t, x, parse.ObjectCloseHint)

	expect.Hint(t, x, parse.ArrayCloseHint)
}
