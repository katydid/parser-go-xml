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
	"fmt"
	"io"

	"github.com/katydid/parser-go-xml/xml/scan"
	"github.com/katydid/parser-go-xml/xml/token"
	"github.com/katydid/parser-go/parse"
)

type Parser interface {
	parse.Parser

	Init(buf []byte)
}

type parser struct {
	// state
	state state
	stack []state

	scanKind  scan.Kind
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
	fmt.Printf("Next %v %d\n", p.state, len(p.stack))
	switch p.state {
	case startState:
		scanKind, err := p.tokenizer.Next()
		if err != nil {
			return parse.UnknownHint, err
		}
		switch scanKind {
		case scan.UnknownKind:
			return parse.UnknownHint, errUknown
		case scan.StartKind:
			p.down(openState)
			return parse.EnterHint, nil
		case scan.AttrKeyKind:
			p.down(attrState)
			return parse.FieldHint, nil
		case scan.AttrValueKind:
			return parse.ValueHint, p.up()
		case scan.CharKind:
			return parse.ValueHint, nil
		case scan.EndKind:
			return parse.LeaveHint, p.up()
		}
	case openState:
		p.down(startState)
		return parse.FieldHint, nil
	case attrState:
		scanKind, err := p.tokenizer.Next()
		if err != nil {
			return parse.UnknownHint, err
		}
		switch scanKind {
		case scan.AttrValueKind:
			return parse.ValueHint, p.up()
		}
		return parse.UnknownHint, errUknown
	case endState:
		return parse.UnknownHint, io.EOF
	}
	return parse.UnknownHint, nil
}

func (p *parser) Skip() error {
	panic("not implemented")
	switch p.state {
	case startState:

	}
	return nil
}

func (p *parser) Token() (parse.Kind, []byte, error) {
	return p.tokenizer.Token()
}

func (p *parser) down(state state) {
	fmt.Printf("down\n")
	// Append the current state to the stack.
	p.stack = append(p.stack, p.state)
	// Create a new state.
	p.state = state
}

func (p *parser) up() error {
	fmt.Printf("up\n")
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
