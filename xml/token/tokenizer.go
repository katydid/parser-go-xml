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

package token

import (
	"strconv"
	"strings"

	"github.com/katydid/parser-go-xml/xml/scan"
	"github.com/katydid/parser-go/cast"
	"github.com/katydid/parser-go/parse"
)

// Tokenizer is a scanner that provides the ability to buffers returned by the scanner into native Go types.
type Tokenizer interface {
	// Next returns the Kind of the token or an error.
	Next() (scan.Kind, error)
	// Token parses and returns the current token.
	Token() (parse.Kind, []byte, error)
}

type tokenizer struct {
	scanner scan.Scanner
	alloc   func(size int) []byte

	scanToken []byte
	scanKind  scan.Kind

	tokenized   bool
	tokenKind   parse.Kind
	tokenErr    error
	tokenDouble float64
	tokenInt    int64
	tokenBytes  []byte
}

func NewTokenizer(buf []byte) Tokenizer {
	alloc := func(size int) []byte {
		return make([]byte, size)
	}
	return &tokenizer{
		scanner: scan.NewScanner(buf),
		alloc:   alloc,
	}
}

// Next returns the Kind of the token or an error.
func (t *tokenizer) Next() (scan.Kind, error) {
	t.tokenized = false
	kind, token, err := t.scanner.Next()
	if err != nil {
		return scan.UnknownKind, err
	}
	t.scanKind = kind
	t.scanToken = []byte(token)
	return kind, nil
}

func (t *tokenizer) Token() (parse.Kind, []byte, error) {
	if err := t.tokenize(); err != nil {
		return parse.UnknownKind, nil, err
	}
	if t.tokenKind == parse.Int64Kind {
		return t.tokenKind, cast.FromInt64(t.tokenInt, t.alloc), nil
	}
	if t.tokenKind == parse.Float64Kind {
		return t.tokenKind, cast.FromFloat64(t.tokenDouble, t.alloc), nil
	}
	return t.tokenKind, t.tokenBytes, nil
}

func (t *tokenizer) tokenize() error {
	if t.scanKind == scan.AttrKeyKind || t.scanKind == scan.StartKind || t.scanKind == scan.EndKind {
		t.tokenKind = parse.StringKind
		t.tokenized = true
		t.tokenBytes = t.scanToken
		return nil
	}
	if t.scanKind == scan.AttrValueKind || t.scanKind == scan.CharKind {
		b, err := strconv.ParseBool(strings.TrimSpace(string(t.scanToken)))
		if err == nil {
			if b {
				t.tokenKind = parse.TrueKind
			} else {
				t.tokenKind = parse.FalseKind
			}
			t.tokenized = true
			return nil
		}
		i, err := strconv.ParseInt(string(t.scanToken), 10, 64)
		if err == nil {
			t.tokenKind = parse.Int64Kind
			t.tokenInt = i
			t.tokenized = true
			return nil
		}
		f, err := strconv.ParseFloat(string(t.scanToken), 64)
		if err == nil {
			t.tokenKind = parse.Float64Kind
			t.tokenDouble = f
			t.tokenized = true
			return nil
		}
		t.tokenKind = parse.StringKind
		t.tokenized = true
		t.tokenBytes = t.scanToken
		return nil
	}
	t.tokenKind = parse.UnknownKind
	t.tokenized = true
	t.tokenErr = ErrNotValue
	return t.tokenErr
}

func (t *tokenizer) Tokenize() (parse.Kind, error) {
	if err := t.tokenize(); err != nil {
		return parse.UnknownKind, err
	}
	return t.tokenKind, nil
}

func (t *tokenizer) Int() (int64, error) {
	if t.tokenKind == parse.Int64Kind {
		return t.tokenInt, nil
	}
	return 0, ErrNotInt
}

func (t *tokenizer) Double() (float64, error) {
	if t.tokenKind == parse.Float64Kind {
		return t.tokenDouble, nil
	}
	return 0, ErrNotDouble
}

func (t *tokenizer) Bytes() ([]byte, error) {
	if t.tokenKind == parse.BytesKind || t.tokenKind == parse.StringKind || t.tokenKind == parse.DecimalKind {
		return t.tokenBytes, nil
	}
	return nil, ErrNotBytes
}
