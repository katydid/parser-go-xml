//  Copyright 2026 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package xml

import (
	"katydid.org.za/go/parser-go-xml/xml/parse"
	goparse "katydid.org.za/go/parser-go/parse"
	"katydid.org.za/go/parser-go/pool"
)

type Parser interface {
	goparse.Parser
	Reset()
	// Init restarts the parser with a new byte buffer, without allocating a new parser.
	Init([]byte)
}

type parser struct {
	parse.Parser
	pool     pool.Pool
	original []byte
}

func (p *parser) Init(buf []byte) {
	p.original = buf
	p.pool.FreeAll()
	p.Parser.Init(buf)
}

func (p *parser) Reset() {
	p.pool.FreeAll()
	p.Parser.Init(p.original)
}

// NewParser returns an XML pull-based parser that skips all spaces.
func NewParser() Parser {
	p := pool.New()
	return &parser{parse.NewParser(parse.WithSkipSpace(), parse.WithAllocator(p.Alloc)), p, nil}
}

// NewRAWParser returns an XML pull-based parser that includes spaces.
func NewRAWParser() Parser {
	p := pool.New()
	return &parser{parse.NewParser(parse.WithAllocator(p.Alloc)), p, nil}
}
