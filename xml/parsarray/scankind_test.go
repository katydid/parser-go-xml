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
	expectHint(t, x, ArrayOpenHint)
	expectHint(t, x, ObjectOpenHint)

	expectHint(t, x, KeyHint)
	expectScanKind(t, x, scan.StartKind)
	expectString(t, x, "a")
	expectHint(t, x, ArrayOpenHint)

	expectHint(t, x, ObjectOpenHint)
	expectHint(t, x, KeyHint)
	expectScanKind(t, x, scan.AttrKeyKind)
	expectString(t, x, "k1")
	expectHint(t, x, ValueHint)
	expectScanKind(t, x, scan.AttrValueKind)
	expectString(t, x, "v1")
	expectHint(t, x, ObjectCloseHint)

	expectHint(t, x, ObjectOpenHint)
	expectHint(t, x, KeyHint)
	expectScanKind(t, x, scan.AttrKeyKind)
	expectString(t, x, "k2")
	expectHint(t, x, ValueHint)
	expectScanKind(t, x, scan.AttrValueKind)
	expectString(t, x, "v2")
	expectHint(t, x, ObjectCloseHint)

	expectHint(t, x, ValueHint)
	expectScanKind(t, x, scan.CharKind)
	expectString(t, x, "b")

	expectHint(t, x, ArrayCloseHint)

	expectHint(t, x, ObjectCloseHint)

	expectHint(t, x, ArrayCloseHint)
}
