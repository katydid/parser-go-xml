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

	xmlparse "github.com/katydid/parser-go-xml/xml/parsarray"
	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
)

func TestTaggerWithElem(t *testing.T) {
	// <a/> will parse as [{"elem": {"a": []}}]
	str := `<a/>`
	x := NewTagger(xmlparse.NewParser(xmlparse.WithBuffer([]byte(str))), WithElemTag())
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.Tag(t, x, "elem")

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "a")
	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestTaggerWithAttr(t *testing.T) {
	// <a k="v"/> will parse as [{"a": [{"attr": {"k": "v"}}]}]
	str := `<a k="v"/>`
	x := NewTagger(xmlparse.NewParser(xmlparse.WithBuffer([]byte(str))), WithAttrTag())
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "a")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.Tag(t, x, "attr")

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "k")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "v")
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestTaggerWithTextElem(t *testing.T) {
	// <a>t</a> will parse as [{"a": [{"text": "t"}]}]
	str := `<a>t</a>`
	x := NewTagger(xmlparse.NewParser(xmlparse.WithBuffer([]byte(str))), WithTextTag())
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "a")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.Tag(t, x, "text")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "t")
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestTaggerWithTextAttr(t *testing.T) {
	// <a k="v"/> will parse as [{"a": [{"k": {"text": "v"}}]}]
	str := `<a k="v"/>`
	x := NewTagger(xmlparse.NewParser(xmlparse.WithBuffer([]byte(str))), WithTextTag())
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "a")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "k")

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.Tag(t, x, "text")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "v")
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestTaggerWithAll(t *testing.T) {
	str := `<a k="v">b</a>`
	// [{"elem": {"a": [{"attr": {"k": {"text": "v"}}}, {"text": "b"}]}}]
	x := NewTagger(xmlparse.NewParser(xmlparse.WithBuffer([]byte(str))), WithAttrTag(), WithElemTag(), WithTextTag())
	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.Tag(t, x, "elem")

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "a")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.Tag(t, x, "attr")
	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "k")
	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.Tag(t, x, "text")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "v")
	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.Tag(t, x, "text")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "b")
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}
