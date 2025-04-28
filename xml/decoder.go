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

package xml

import (
	"bytes"
	"encoding/xml"
	"io"
)

type decoder interface {
	Skip() error
	Token() (xml.Token, error)
}

type xmlDecoder struct {
	*xml.Decoder
}

func newDecoder(buf []byte) decoder {
	buf = removeProcessingInstructions(buf)
	buf = bytes.TrimSpace(buf)
	dec := xml.NewDecoder(bytes.NewBuffer(buf))
	dec.Strict = false
	return &xmlDecoder{dec}
}

// Token skips all the uninteresting tokens and only returns StartElement, non empty CharData's or EndElement
func (x *xmlDecoder) Token() (xml.Token, error) {
	return scan(x.Decoder)
}

func hasContent(c xml.CharData) bool {
	return len(string(c)) > 0
}

func scan(dec *xml.Decoder) (xml.Token, error) {
	tok, err := dec.Token()
	for {
		if err != nil {
			return tok, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			return tok, nil
		case xml.CharData:
			if hasContent(t) {
				return tok, nil
			}
		case xml.EndElement:
			return tok, io.EOF
		}
		tok, err = dec.Token()
	}
}
