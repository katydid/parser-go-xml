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

package tag

import (
	xmlparse "github.com/katydid/parser-go-xml/xml/parse"
	"github.com/katydid/parser-go/parse"
)

type tagger struct {
	parser  xmlparse.Parser
	options *options
}

func NewTagger(p xmlparse.Parser, opts ...Option) parse.Parser {
	options := newOptions(opts...)
	return &tagger{
		parser:  p,
		options: options,
	}
}

func (t *tagger) Next() (parse.Hint, error) {
	return t.parser.Next()
}

func (t *tagger) Skip() error {
	return t.parser.Skip()
}

func (t *tagger) Token() (parse.Kind, []byte, error) {
	kind, val, err := t.parser.Token()
	if err != nil {
		return kind, val, err
	}
	switch t.parser.TokenXMLType() {
	case xmlparse.UnknownXMLType:
		return kind, val, nil
	case xmlparse.TextXMLType:
		if len(t.options.textPrefix) == 0 {
			return kind, val, nil
		}
		return kind, copyAppend(t.options.textPrefix, val), nil
	case xmlparse.ElemXMLType:
		if len(t.options.elemPrefix) == 0 {
			return kind, val, nil
		}
		return kind, copyAppend(t.options.elemPrefix, val), nil
	case xmlparse.AttrXMLType:
		if len(t.options.attrPrefix) == 0 {
			return kind, val, nil
		}
		return kind, copyAppend(t.options.attrPrefix, val), nil
	}
	panic("unreachable")
}

func copyAppend(b1, b2 []byte) []byte {
	res := make([]byte, 0, len(b1)+len(b2))
	res = append(res, b1...)
	res = append(res, b2...)
	return res
}
