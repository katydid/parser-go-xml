//  Copyright 2015 Walter Schulze
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

// Package xml contains a parser for XML.
package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"katydid.org.za/go/parser-go/parser"
)

type xmlParser struct {
	dec decoder

	// options
	attrPrefix string
	elemPrefix string
	textPrefix string

	tok       xml.Token
	attrs     []xml.Attr
	attrIndex int
	attrValue bool
	attrFirst bool
}

// XMLParser is an xml parser.
type XMLParser interface {
	parser.Interface
	//Init intialises the parser with a byte buffer containing xml.
	Init([]byte) error
}

// NewXMLParser returns a new xml parser.
func NewXMLParser(opts ...Option) XMLParser {
	options := newOptions(opts...)
	x := &xmlParser{}
	x.attrPrefix = options.attrPrefix
	x.elemPrefix = options.elemPrefix
	x.textPrefix = options.textPrefix
	return x
}

func (p *xmlParser) Init(buf []byte) error {
	p.dec = newDecoder(buf)
	return nil
}

func (p *xmlParser) Next() (err error) {
	if p.attrValue {
		if p.attrFirst {
			p.attrFirst = false
			return nil
		} else {
			return io.EOF
		}
	}
	if p.tok == nil && p.attrs != nil {
		p.attrIndex++
		if p.attrIndex < len(p.attrs) {
			return nil
		}
	}
	if p.tok != nil {
		for {
			if _, ok := p.tok.(xml.StartElement); ok {
				if err := p.dec.Skip(); err != nil {
					return err
				}
				break
			} else if c, ok := p.tok.(xml.CharData); ok {
				if hasContent(c) {
					break
				}
			} else {
				panic(fmt.Sprintf("unknown token %T", p.tok))
			}
		}
	}
	p.tok, err = p.dec.Token()
	return err
}

func (p *xmlParser) IsLeaf() bool {
	if p.tok == nil {
		return p.attrValue
	}
	_, ok := p.tok.(xml.CharData)
	return ok
}

func (p *xmlParser) getValue() string {
	if p.tok == nil && p.attrValue {
		return p.attrs[p.attrIndex].Value
	}
	if c, ok := p.tok.(xml.CharData); ok {
		return string(c)
	}
	return ""
}

func (p *xmlParser) Double() (float64, error) {
	return strconv.ParseFloat(p.getValue(), 64)
}

func (p *xmlParser) Int() (int64, error) {
	i, err := strconv.ParseInt(p.getValue(), 10, 64)
	return int64(i), err
}

func (p *xmlParser) Uint() (uint64, error) {
	i, err := strconv.ParseUint(p.getValue(), 10, 64)
	return uint64(i), err
}

func (p *xmlParser) Bool() (bool, error) {
	return strconv.ParseBool(strings.TrimSpace(p.getValue()))
}

func (p *xmlParser) String() (string, error) {
	if p.tok == nil && p.attrIndex < len(p.attrs) {
		if p.attrValue {
			return p.textPrefix + p.attrs[p.attrIndex].Value, nil
		} else {
			return p.attrPrefix + p.attrs[p.attrIndex].Name.Local, nil
		}
	}
	if s, ok := p.tok.(xml.StartElement); ok {
		return p.elemPrefix + s.Name.Local, nil
	}
	if c, ok := p.tok.(xml.CharData); ok {
		return p.textPrefix + string(c), nil
	}
	return "", parser.ErrNotString
}

func (p *xmlParser) Bytes() ([]byte, error) {
	if c, ok := p.tok.(xml.CharData); ok {
		return []byte(c), nil
	}
	return nil, parser.ErrNotBytes
}

func (p *xmlParser) Up() {
	if p.tok == nil {
		if p.attrValue {
			p.attrValue = false
			p.attrFirst = false
			return
		}
	}
	if _, ok := p.tok.(xml.EndElement); ok {
		p.tok = nil
		p.attrs = nil
		p.attrIndex = 0
		return
	}
	if err := p.dec.Skip(); err != nil {
		if err != io.EOF {
			panic(err)
		}
	}
}

func (p *xmlParser) Down() {
	if p.tok == nil {
		if p.attrIndex < len(p.attrs) {
			p.attrValue = true
			p.attrFirst = true
			return
		}
	}
	if s, ok := p.tok.(xml.StartElement); ok {
		p.tok = nil
		p.attrs = s.Attr
		p.attrIndex = -1
		return
	}
	panic(fmt.Sprintf("not a start element %T", p.tok))
}
