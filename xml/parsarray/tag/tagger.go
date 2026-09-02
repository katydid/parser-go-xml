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
	"fmt"
	"io"

	xmlparse "github.com/katydid/parser-go-xml/xml/parsarray"
	"github.com/katydid/parser-go-xml/xml/scan"
	"katydid.org.za/go/parser-go/parse"
)

type tagger struct {
	parser  xmlparse.Parser
	options *options

	// state
	state state
	stack []state
}

type Tagger interface {
	parse.Parser
	Init(buf []byte)
}

func NewTagger(p xmlparse.Parser, opts ...Option) Tagger {
	options := newOptions(opts...)
	return &tagger{
		parser:  p,
		options: options,
		state:   startState,
		stack:   make([]state, 0, 10),
	}
}

func (t *tagger) Init(buf []byte) {
	t.state = startState
	t.parser.Init(buf)
	t.stack = t.stack[:0]
}

func (t *tagger) Next() (parse.Hint, error) {
	switch t.state {
	case startState:
		h, err := t.parser.Next()
		if err != nil {
			return parse.UnknownHint, err
		}
		switch t.parser.ScanKind() {
		case scan.UnknownKind:
			return translateHint(h), nil
		case scan.StartKind:
			if h != xmlparse.ObjectOpenHint {
				return translateHint(h), nil
			}
			if !t.options.tagElems {
				return translateHint(h), nil
			}
			t.down(elemTagOpenState)
			return parse.EnterHint, nil
		case scan.EndKind:
			if h != xmlparse.ObjectCloseHint {
				return translateHint(h), nil
			}
			if !t.options.tagElems {
				return translateHint(h), nil
			}
			t.state = elemTagCloseState
			return parse.LeaveHint, nil
		case scan.AttrKeyKind:
			if h != xmlparse.ObjectOpenHint {
				return translateHint(h), nil
			}
			if !t.options.tagAttrs {
				return translateHint(h), nil
			}
			t.down(attrTagOpenState)
			return parse.EnterHint, nil
		case scan.AttrValueKind:
			if h == xmlparse.ObjectCloseHint {
				if !t.options.tagAttrs {
					return translateHint(h), nil
				}
				t.state = attrTagCloseState
				return parse.LeaveHint, nil
			} else if h == xmlparse.ValueHint {
				if !t.options.tagTexts {
					return translateHint(h), nil
				}
				t.state = textTagOpenState
				return parse.EnterHint, nil
			}
			return translateHint(h), nil
		case scan.CharKind:
			if h != xmlparse.ValueHint {
				return translateHint(h), nil
			}
			if !t.options.tagTexts {
				return translateHint(h), nil
			}
			t.state = textTagOpenState
			return parse.EnterHint, nil
		}
		panic("unreachable")
	case elemTagOpenState:
		t.state = elemTagKeyState
		return parse.FieldHint, nil
	case elemTagKeyState:
		t.state = startState
		return parse.EnterHint, nil
	case elemTagCloseState:
		t.up()
		return parse.LeaveHint, nil
	case attrTagOpenState:
		t.state = attrTagKeyState
		return parse.FieldHint, nil
	case attrTagKeyState:
		t.state = startState
		return parse.EnterHint, nil
	case attrTagCloseState:
		t.up()
		return parse.LeaveHint, nil
	case textTagOpenState:
		t.state = textTagKeyState
		return parse.FieldHint, nil
	case textTagKeyState:
		t.state = textTagCloseState
		return parse.ValueHint, nil
	case textTagCloseState:
		t.state = startState
		return parse.LeaveHint, nil
	case endState:
		return parse.UnknownHint, io.EOF
	}
	panic(fmt.Sprintf("unreachable: unknown state = %c", t.state))
}

func (t *tagger) Skip() error {
	return t.parser.Skip()
}

func (t *tagger) Token() (parse.Kind, []byte, error) {
	switch t.state {
	case elemTagKeyState:
		return parse.TagKind, []byte("elem"), nil
	case attrTagKeyState:
		return parse.TagKind, []byte("attr"), nil
	case textTagKeyState:
		return parse.TagKind, []byte("text"), nil
	}
	kind, val, err := t.parser.Token()
	if err != nil {
		return kind, val, err
	}
	switch t.parser.ScanKind() {
	case scan.UnknownKind:
		return kind, val, err
	case scan.StartKind:
		if len(t.options.elemPrefix) == 0 {
			return kind, val, nil
		}
		return kind, copyAppend(t.options.elemPrefix, val), nil
	case scan.EndKind:
		return kind, val, nil
	case scan.AttrKeyKind:
		if len(t.options.attrPrefix) == 0 {
			return kind, val, nil
		}
		return kind, copyAppend(t.options.attrPrefix, val), nil
	case scan.AttrValueKind:
		if len(t.options.textPrefix) == 0 {
			return kind, val, nil
		}
		return kind, copyAppend(t.options.textPrefix, val), nil
	case scan.CharKind:
		if len(t.options.textPrefix) == 0 {
			return kind, val, nil
		}
		return kind, copyAppend(t.options.textPrefix, val), nil
	}
	panic("unreachable")
}

func (t *tagger) down(state state) {
	// Append the current state to the stack.
	t.stack = append(t.stack, t.state)
	// Create a new state.
	t.state = state
}

func (t *tagger) up() error {
	top := len(t.stack) - 1
	// Set the current state to the state on top of the stack.
	t.state = t.stack[top]
	// Remove the state on the top the stack from the stack,
	// but do it in a way that keeps the capacity,
	// so we can reuse it the next time Down is called.
	t.stack = t.stack[:top]
	return nil
}

func copyAppend(b1, b2 []byte) []byte {
	res := make([]byte, 0, len(b1)+len(b2))
	res = append(res, b1...)
	res = append(res, b2...)
	return res
}
