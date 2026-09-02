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

package tag

import (
	"testing"

	xmlparse "katydid.org.za/go/parser-go-xml/xml/parsarray"
	"katydid.org.za/go/parser-go/expect"
	"katydid.org.za/go/parser-go/parse"
)

func TestTaggerPrefix(t *testing.T) {
	str := `<a k1="v1" k2="v2">b</a>`
	x := NewTagger(xmlparse.NewParser(xmlparse.WithBuffer([]byte(str))), WithAttrPrefix("attr_"), WithElemPrefix("elem_"), WithTextPrefix("text_"))
	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "elem_a")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "attr_k1")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "text_v1")
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "attr_k2")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "text_v2")
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "text_b")

	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)
}
