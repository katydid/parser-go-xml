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
	Next() (Token, error)
}

type scanner struct {
	d       *xml.Decoder
	current xml.Token
	attrs   *attrs
}

func NewScanner(buf []byte) Scanner {
	buf = removeProcessingInstructions(buf)
	buf = bytes.TrimSpace(buf)
	dec := xml.NewDecoder(bytes.NewBuffer(buf))
	dec.Strict = false
	return &scanner{d: dec}
}

func (s *scanner) Next() (Token, error) {
	if s.current == nil {
		if err := s.next(); err != nil {
			return Token{Typ: UnknownToken}, err
		}
	}
	switch t := s.current.(type) {
	case xml.StartElement:
		if s.attrs == nil {
			s.attrs = newAttrs(t.Attr)
			return Token{
				Typ: StartToken,
				Val: t.Name.Local,
			}, nil
		}
		tok, err := s.attrs.Next()
		if err == nil {
			return tok, nil
		}
		if err != io.EOF {
			return Token{Typ: UnknownToken}, err
		}
		s.reset()
		return s.Next()
	case xml.CharData:
		token := Token{
			Typ: CharToken,
			Val: string(t),
		}
		s.reset()
		return token, nil
	case xml.EndElement:
		token := Token{
			Typ: EndToken,
			Val: t.Name.Local,
		}
		s.reset()
		return token, nil
	}
	return Token{Typ: UnknownToken}, nil
}

func (s *scanner) reset() {
	s.current = nil
	s.attrs = nil
}

// next only returns a StartElement, non empty CharData or an EndElement
func (s *scanner) next() error {
	tok, err := s.d.Token()
	for {
		if err != nil {
			s.current = &Token{Typ: UnknownToken}
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			s.current = tok
			return nil
		case xml.CharData:
			if hasContent(t) {
				s.current = tok
				return nil
			}
		case xml.EndElement:
			s.current = tok
			return nil
		}
		tok, err = s.d.Token()
	}
}

func hasContent(c xml.CharData) bool {
	return len(string(c)) > 0
}
