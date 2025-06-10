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
	"errors"
	"io"

	"github.com/katydid/parser-go-xml/xml/scan"
	"github.com/katydid/parser-go-xml/xml/token"
	"github.com/katydid/parser-go/parse"
)

type Parser interface {
	parse.Parser
	// Internal: only for internal use
	ScanKind() scan.Kind
}

type parser struct {
	// state
	state state
	stack []state

	scanKind  scan.Kind
	tokenizer token.Tokenizer
}

func NewParser(buf []byte) Parser {
	p := &parser{
		state:     startState,
		stack:     make([]state, 0, 10),
		tokenizer: token.NewTokenizer(buf),
	}
	return p
}

func (p *parser) Next() (parse.Hint, error) {
	switch p.state {
	case startState:
		return p.nextStart()
	case startedState:
		return p.nextStarted()
	case arrayOpenedState:
		return p.nextArrayOpened()
	case objectOpenedState:
		return p.nextObjectOpened()
	case objectKeyedState:
		return p.nextObjectKeyed()
	case objectValuedState:
		return p.nextObjectValued()
	case attrOpenedState:
		return p.nextAttrOpened()
	case attrKeyedState:
		return p.nextAttrKeyed()
	case attrValuedState:
		return p.nextAttrValued()
	case endState:
		return parse.UnknownHint, io.EOF
	default:
		panic("unreachable")
	}
}

func (p *parser) nextStart() (parse.Hint, error) {
	p.down(startedState)
	return parse.ArrayOpenHint, nil
}

func (p *parser) nextStarted() (parse.Hint, error) {
	scanKind, err := p.tokenizer.Next()
	p.scanKind = scanKind
	if err != nil {
		if err == io.EOF {
			if err := p.up(); err != nil {
				return parse.UnknownHint, err
			}
			return parse.ArrayCloseHint, nil
		}
		return parse.UnknownHint, err
	}
	switch scanKind {
	case scan.StartKind:
		p.down(objectOpenedState)
		return parse.ObjectOpenHint, nil
	case scan.CharKind:
		return parse.ValueHint, nil
	case scan.AttrKeyKind, scan.AttrValueKind, scan.EndKind:
		return parse.UnknownHint, errors.New("expected start element or characters")
	default:
		panic("unreachable")
	}
}

func (p *parser) nextArrayOpened() (parse.Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return parse.UnknownHint, err
	}
	switch scanKind {
	case scan.StartKind:
		p.down(objectOpenedState)
		return parse.ObjectOpenHint, nil
	case scan.CharKind:
		return parse.ValueHint, nil
	case scan.AttrKeyKind:
		p.down(attrOpenedState)
		return parse.ObjectOpenHint, nil
	case scan.AttrValueKind:
		return parse.UnknownHint, errors.New("expected start element, characters, attribute key or end element")
	case scan.EndKind:
		if err := p.up(); err != nil {
			return parse.UnknownHint, err
		}
		return parse.ArrayCloseHint, nil
	default:
		panic("unreachable")
	}
}

func (p *parser) nextAttrOpened() (parse.Hint, error) {
	p.state = attrKeyedState
	return parse.KeyHint, nil
}

func (p *parser) nextAttrKeyed() (parse.Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return parse.UnknownHint, err
	}
	if scanKind != scan.AttrValueKind {
		return parse.UnknownHint, errors.New("expected attr value")
	}
	p.state = attrValuedState
	return parse.ValueHint, nil
}

func (p *parser) nextAttrValued() (parse.Hint, error) {
	if err := p.up(); err != nil {
		return parse.UnknownHint, err
	}
	return parse.ObjectCloseHint, nil
}

func (p *parser) nextObjectOpened() (parse.Hint, error) {
	p.state = objectKeyedState
	return parse.KeyHint, nil
}

func (p *parser) nextObjectKeyed() (parse.Hint, error) {
	p.state = objectValuedState
	p.down(arrayOpenedState)
	return parse.ArrayOpenHint, nil
}

func (p *parser) nextObjectValued() (parse.Hint, error) {
	if err := p.up(); err != nil {
		return parse.UnknownHint, err
	}
	return parse.ObjectCloseHint, nil
}

func (p *parser) Skip() error {
	switch p.state {
	case startedState, arrayOpenedState:
		// '[' has been parsed or
		// '['"e1",...,"en" has been parsed.
		// call Next until ']' is parsed,
		// which will result in the stack being popped,
		// which will result in the stack size being smaller.
		currentStackSize := len(p.stack)
		for len(p.stack) >= currentStackSize {
			_, err := p.Next()
			if err != nil {
				return err
			}
		}
	case objectOpenedState, objectValuedState, attrOpenedState, attrValuedState:
		// '{' has been parsed or
		// '{'"k1":"v1",...,"kn":"vn" has been parsed.
		// call Next until '}' is parsed,
		// which will result in the stack being popped,
		// which will result in the stack size being smaller.
		currentStackSize := len(p.stack)
		for len(p.stack) >= currentStackSize {
			_, err := p.Next()
			if err != nil {
				return err
			}
		}
	case objectKeyedState, attrKeyedState:
		currentStackSize := len(p.stack)
		_, err := p.Next()
		if err != nil {
			return err
		}
		// If Next parsed down into an array or object,
		// then keep on parsing until we reach our current level.
		// If Next parsed a string, number, boolean or null,
		// then the level would be the same.
		for len(p.stack) > currentStackSize {
			_, err := p.Next()
			if err != nil {
				return err
			}
		}
	default:
		_, err := p.Next()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) Token() (parse.Kind, []byte, error) {
	return p.tokenizer.Token()
}

// Internal: only for internal use
// ScanKind is only used by the Tagger,
// which needs some internal knowledge to properly apply tags.
func (p *parser) ScanKind() scan.Kind {
	return p.scanKind
}

func (p *parser) nextToken() (scan.Kind, error) {
	scanKind, err := p.tokenizer.Next()
	p.scanKind = scanKind
	if err == nil {
		return scanKind, nil
	}
	if err == io.EOF {
		return scanKind, io.ErrShortBuffer
	}
	return scanKind, err
}

func (p *parser) down(state state) {
	// Append the current state to the stack.
	p.stack = append(p.stack, p.state)
	// Create a new state.
	p.state = state
}

func (p *parser) up() error {
	top := len(p.stack) - 1
	// Set the current state to the state on top of the stack.
	p.state = p.stack[top]
	// Remove the state on the top the stack from the stack,
	// but do it in a way that keeps the capacity,
	// so we can reuse it the next time Down is called.
	p.stack = p.stack[:top]
	if len(p.stack) == 0 {
		p.state = endState
	}
	return nil
}
