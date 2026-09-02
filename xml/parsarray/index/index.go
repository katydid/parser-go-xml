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
	"fmt"
	"io"

	xmlparse "katydid.org.za/go/parser-go-xml/xml/parsarray"
	"katydid.org.za/go/parser-go/cast"
	"katydid.org.za/go/parser-go/parse"
)

type Parser interface {
	parse.Parser
}

type parser struct {
	parser xmlparse.Parser

	// state
	state state
	stack []state
}

func WithIndexedArrays(p xmlparse.Parser) parse.Parser {
	return &parser{
		parser: p,
		state:  state{},
		stack:  make([]state, 0, 10),
	}
}

func (p *parser) Next() (parse.Hint, error) {
	switch p.state.kind {
	case startState:
		h, err := p.parser.Next()
		if err != nil {
			return parse.UnknownHint, err
		}
		switch h {
		case xmlparse.ObjectOpenHint:
			p.down(startState)
			return parse.EnterHint, nil
		case xmlparse.ObjectCloseHint:
			if err := p.up(); err != nil {
				return parse.UnknownHint, err
			}
			return parse.LeaveHint, nil
		case xmlparse.ArrayOpenHint:
			p.down(arrayIndexState)
			return parse.EnterHint, nil
		case xmlparse.ArrayCloseHint:
			if err := p.up(); err != nil {
				return parse.UnknownHint, err
			}
			return parse.LeaveHint, nil
		}
		return translateHint(h), nil
	case arrayIndexState:
		h, err := p.parser.Next()
		if err != nil {
			return parse.UnknownHint, err
		}
		p.state.hint = h
		if p.state.hint == xmlparse.ArrayCloseHint {
			if err := p.up(); err != nil {
				return parse.UnknownHint, err
			}
			return parse.LeaveHint, nil
		}
		p.state.index++
		p.state.kind = arrayElemState
		return parse.FieldHint, nil
	case arrayElemState:
		p.state.kind = arrayIndexState
		h := p.state.hint
		switch h {
		case xmlparse.ObjectOpenHint:
			p.down(startState)
			return parse.EnterHint, nil
		case xmlparse.ObjectCloseHint:
			if err := p.up(); err != nil {
				return parse.UnknownHint, err
			}
			return parse.LeaveHint, nil
		case xmlparse.ArrayOpenHint:
			p.down(arrayIndexState)
			return parse.EnterHint, nil
		case xmlparse.ArrayCloseHint:
			if err := p.up(); err != nil {
				return parse.UnknownHint, err
			}
			return parse.LeaveHint, nil
		}
		return translateHint(h), nil
	case endState:
		return parse.UnknownHint, io.EOF
	}
	panic(fmt.Sprintf("unreachable: unknown state = %v", p.state))
}

func (p *parser) Skip() error {
	switch p.state.kind {
	case startState:
		if len(p.stack) == 0 {
			_, err := p.Next()
			return err
		}
		if p.state.hint != xmlparse.KeyHint {
			// do not go up when it is an object value that needs to be skipped over
			if err := p.up(); err != nil {
				return err
			}
		}
		p.state.hint = xmlparse.UnknownHint
		return p.parser.Skip()
	case arrayIndexState:
		if err := p.up(); err != nil {
			return err
		}
		return p.parser.Skip()
	case arrayElemState:
		p.state.kind = arrayIndexState
		if p.state.hint == xmlparse.ValueHint {
			// values do not need to be skipped, Next will take care of it.
			return nil
		}
		return p.parser.Skip()
	case endState:
		return p.parser.Skip()
	}
	panic(fmt.Sprintf("unreachable: unknown state = %v", p.state))
}

func alloc(size int) []byte { return make([]byte, size) }

func (p *parser) Token() (parse.Kind, []byte, error) {
	if p.state.kind == arrayElemState {
		return parse.Int64Kind, cast.FromInt64Ptr(&p.state.index, alloc), nil
	}
	return p.parser.Token()
}

func (p *parser) down(stateKind stateKind) {
	// Append the current state to the stack.
	p.stack = append(p.stack, p.state)
	// Create a new state.
	p.state.kind = stateKind
	p.state.index = -1
}

func (p *parser) up() error {
	if len(p.stack) == 0 {
		return errUnexpectedClose
	}
	top := len(p.stack) - 1
	// Set the current state to the state on top of the stack.
	p.state = p.stack[top]
	// Remove the state on the top the stack from the stack,
	// but do it in a way that keeps the capacity,
	// so we can reuse it the next time down is called.
	p.stack = p.stack[:top]
	if len(p.stack) == 0 {
		p.state.kind = endState
	}
	return nil
}
