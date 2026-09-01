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

package scan

import (
	"bytes"
	"encoding/xml"
	"io"
)

type Scanner interface {
	Next() (Kind, string, error)
	Peek() (Kind, error)
}

type ScannerWithInit interface {
	Scanner
	Init([]byte)
}

type scanner struct {
	*options
	dec     *xml.Decoder
	current xml.Token
	attrs   *attrs
	peek    xml.Token
	peekerr error
}

func NewScanner(opts ...Option) ScannerWithInit {
	options := newOptions(opts...)
	s := &scanner{options: options}
	if options.buf != nil {
		s.Init(options.buf)
	}
	return s
}

func (s *scanner) Init(buf []byte) {
	buf = removeProcessingInstructions(buf)
	buf = bytes.TrimSpace(buf)
	s.dec = xml.NewDecoder(bytes.NewBuffer(buf))
	s.dec.Strict = false
}

func (s *scanner) Next() (Kind, string, error) {
	if s.current == nil {
		var err error
		s.current, err = s.token()
		if err != nil {
			return UnknownKind, "", err
		}
	}
	switch t := s.current.(type) {
	case xml.StartElement:
		if s.attrs == nil {
			s.attrs = newAttrs(t.Attr)
			return StartKind, t.Name.Local, nil
		}
		kind, val, err := s.attrs.Next()
		if err == nil {
			return kind, val, nil
		}
		if err != io.EOF {
			return UnknownKind, "", err
		}
		s.reset()
		return s.Next()
	case xml.CharData:
		s.reset()
		return CharKind, string(t), nil
	case xml.EndElement:
		s.reset()
		return EndKind, t.Name.Local, nil
	}
	return UnknownKind, "", nil
}

func (s *scanner) reset() {
	s.current = nil
	s.attrs = nil
}

func (s *scanner) Peek() (Kind, error) {
	if s.peek == nil && s.peekerr == nil {
		if s.current != nil {
			if s.attrs != nil {
				k, err := s.attrs.Peek()
				if err == nil {
					return k, nil
				}
			}
			s.current = xml.CopyToken(s.current)
		}
		s.peek, s.peekerr = s.token()
	}
	if s.peekerr != nil {
		return UnknownKind, s.peekerr
	}
	switch s.peek.(type) {
	case xml.StartElement:
		return StartKind, nil
	case xml.CharData:
		return CharKind, nil
	case xml.EndElement:
		return EndKind, nil
	}
	return UnknownKind, nil
}

// token returns the next token and takes into account the tokens that were peeked.
// It only sets s.current to a StartElement, non empty CharData or an EndElement.
// It does not look at Attributes.
func (s *scanner) token() (xml.Token, error) {
	if s.peek != nil || s.peekerr != nil {
		tok := s.peek
		err := s.peekerr
		s.peek = nil
		s.peekerr = nil
		return tok, err
	}
	tok, err := s.nextToken()
	return tok, err
}

// nextToken asks the decoder to decode the next token.
func (s *scanner) nextToken() (xml.Token, error) {
	tok, err := s.dec.Token()
	if err != nil {
		return nil, err
	}
	if c, ok := tok.(xml.CharData); ok {
		if len(string(c)) == 0 {
			return s.nextToken()
		}
		if s.skipSpace && len(bytes.TrimSpace(c)) == 0 {
			return s.nextToken()
		}
	}
	return tok, nil
}
