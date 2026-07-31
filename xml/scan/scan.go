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
}

type ScannerWithInit interface {
	Scanner
	Init([]byte)
}

type scanner struct {
	*options
	d       *xml.Decoder
	current xml.Token
	attrs   *attrs
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
	dec := xml.NewDecoder(bytes.NewBuffer(buf))
	dec.Strict = false
	s.d = dec
}

func (s *scanner) Peek() (Kind, error) {
	if s.current == nil {
		if err := s.next(); err != nil {
			return UnknownKind, err
		}
	}
	switch t := s.current.(type) {
	case xml.StartElement:
		if s.attrs == nil {
			s.attrs = newAttrs(t.Attr)
			return StartKind, nil
		}
		kind, err := s.attrs.Peak()
		if err == nil {
			return kind, nil
		}
		if err != io.EOF {
			return UnknownKind, err
		}
		s.reset()
		return s.Peek()
	case xml.CharData:
		if s.skipSpace && len(bytes.TrimSpace(t)) == 0 {
			s.reset()
			return s.Peek()
		}
		s.reset()
		return CharKind, nil
	case xml.EndElement:
		s.reset()
		return EndKind, nil
	}
	return UnknownKind, nil
}

func (s *scanner) Next() (Kind, string, error) {
	if s.current == nil {
		if err := s.next(); err != nil {
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
		if s.skipSpace && len(bytes.TrimSpace(t)) == 0 {
			s.reset()
			return s.Next()
		}
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

// next only returns a StartElement, non empty CharData or an EndElement
func (s *scanner) next() error {
	tok, err := s.d.Token()
	for {
		if err != nil {
			s.current = tok
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
