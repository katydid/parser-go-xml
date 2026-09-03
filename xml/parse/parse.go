// Copyright 2026 Walter Schulze
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

	"katydid.org.za/go/parser-go-xml/xml/scan"
	"katydid.org.za/go/parser-go-xml/xml/token"
	"katydid.org.za/go/parser-go/parse"
)

type Parser interface {
	parse.Parser

	Init(buf []byte)
}

type parser struct {
	state state
	stack []state

	peekKind  scan.Kind
	peekErr   error
	tokenizer token.TokenizerWithInit
}

func NewParser(opts ...Option) Parser {
	options := newOptions(opts...)
	p := &parser{
		stack:     make([]state, 0, 10),
		tokenizer: token.NewTokenizer(options.toTokenOptions()...),
	}
	if options.buf != nil {
		p.Init(options.buf)
	}
	return p
}

func (p *parser) Init(buf []byte) {
	p.stack = p.stack[:0]
	p.state = startState
	p.tokenizer.Init(buf)
}

func (p *parser) Next() (parse.Hint, error) {
	switch p.state {
	case startState:
		p.state = endState
		p.down(inElemState)
		return parse.EnterHint, nil
	case inElemState:
		scanKind, err := p.next()
		if err == io.EOF {
			if err := p.up(); err != nil {
				return parse.UnknownHint, err
			}
			return parse.LeaveHint, nil
		}
		if err != nil {
			return parse.UnknownHint, err
		}
		switch scanKind {
		case scan.UnknownKind:
			return parse.UnknownHint, errUknown
		case scan.StartKind:
			p.state = atFieldState
			return parse.FieldHint, nil
		case scan.AttrKeyKind:
			p.state = atAttributeKeyState
			return parse.FieldHint, nil
		case scan.AttrValueKind:
			return parse.UnknownHint, errors.New("unexpected attr value")
		case scan.CharKind:
			return parse.ValueHint, nil
		case scan.EndKind:
			if err := p.up(); err != nil {
				return parse.UnknownHint, err
			}
			return parse.LeaveHint, nil
		}
		panic("unreachable")
	case atFieldState:
		nextKind1, nextKind2, err := p.look2()
		if err == nil && nextKind1 == scan.CharKind && nextKind2 == scan.EndKind {
			// This is a leaf/value so going to the next token is the appropriate action.
			p.next()
			p.state = inElemState
			p.down(isLeafState)
			return parse.ValueHint, nil
		}
		// This is not an element with another element or multiple values, so we need to enter it.
		p.state = inElemState
		p.down(inElemState)
		return parse.EnterHint, nil
	case isLeafState:
		scanKind, err := p.next()
		if err != nil {
			return parse.UnknownHint, err
		}
		if scanKind != scan.EndKind {
			return parse.UnknownHint, errors.New("expected attr value")
		}
		if err := p.up(); err != nil {
			return parse.UnknownHint, err
		}
		return p.Next()
	case atAttributeKeyState:
		scanKind, err := p.next()
		if err != nil {
			return parse.UnknownHint, err
		}
		if scanKind != scan.AttrValueKind {
			return parse.UnknownHint, errors.New("expected attr value")
		}
		p.state = inElemState
		p.down(atAttributeValueState)
		return parse.ValueHint, nil
	case atAttributeValueState:
		if err := p.up(); err != nil {
			return parse.UnknownHint, err
		}
		return p.Next()
	case endState:
		return parse.UnknownHint, io.EOF
	}
	return parse.UnknownHint, nil
}

func (p *parser) Skip() error {
	switch p.state {
	case startState:
		p.state = endState
		return nil
	case inElemState:
		// <A>... call until </A> is parsed
		// <A><B>C</B>... call until </A> is parsed
		currentStackSize := len(p.stack)
		for len(p.stack) >= currentStackSize {
			_, err := p.Next()
			if err != nil {
				return err
			}
		}
		return nil
	case atFieldState:
		// go down into the field
		_, err := p.Next()
		if err != nil {
			return err
		}
		switch p.state {
		case isLeafState:
		case inElemState:
			return p.Skip()
		default:
			panic("unreachable")
		}
		return nil
	case isLeafState:
		_, err := p.Next()
		if err != nil {
			return err
		}
		currentStackSize := len(p.stack)
		for len(p.stack) >= currentStackSize {
			_, err := p.Next()
			if err != nil {
				return err
			}
		}
		return nil
	case atAttributeKeyState:
		_, err := p.Next()
		return err
	case atAttributeValueState:
		// <A b="c" .. call until </A> is parsed
		_, err := p.Next()
		if err != nil {
			return err
		}
		currentStackSize := len(p.stack)
		for len(p.stack) >= currentStackSize {
			_, err := p.Next()
			if err != nil {
				return err
			}
		}
		return nil
	case endState:
		return nil
	}
	panic("unreachable")
}

func (p *parser) Token() (parse.Kind, []byte, error) {
	return p.tokenizer.Token()
}

// next returns the next scan kind and takes peek into account.
func (p *parser) next() (scan.Kind, error) {
	if p.peekKind != scan.UnknownKind || p.peekErr != nil {
		k, err := p.peekKind, p.peekErr
		p.peekKind, p.peekErr = scan.UnknownKind, nil
		return k, err
	}
	k, err := p.tokenizer.Next()
	return k, err
}

// look2 returns the next two scan kinds
func (p *parser) look2() (scan.Kind, scan.Kind, error) {
	if p.peekKind == scan.UnknownKind && p.peekErr == nil {
		p.peekKind, p.peekErr = p.tokenizer.Next()
		if p.peekErr != nil {
			return scan.UnknownKind, scan.UnknownKind, nil
		}
	}
	peekKind2, peekErr2 := p.tokenizer.Peek()
	return p.peekKind, peekKind2, peekErr2
}

func (p *parser) down(state state) {
	// Append the current state to the stack.
	p.stack = append(p.stack, p.state)
	// Create a new state.
	p.state = state
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
	// so we can reuse it the next time Down is called.
	p.stack = p.stack[:top]
	if len(p.stack) == 0 {
		p.state = endState
	}
	return nil
}
