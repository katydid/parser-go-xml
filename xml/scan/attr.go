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
	"encoding/xml"
	"io"
)

type attrs struct {
	attrs   []*attr
	index   int
	atValue bool
}

func newAttrs(xmlattrs []xml.Attr) *attrs {
	attrlist := make([]*attr, len(xmlattrs))
	for i := range xmlattrs {
		attrlist[i] = newAttr(xmlattrs[i])
	}
	return &attrs{attrs: attrlist, index: 0, atValue: false}
}

func (a *attrs) Next() (Token, error) {
	if a.index >= len(a.attrs) {
		return Token{Typ: UnknownToken}, io.EOF
	}
	if a.atValue {
		// we are at the value, so return the value
		token := Token{
			Typ: AttrValueToken,
			Val: a.attrs[a.index].val,
		}
		// move onto next key
		a.atValue = false
		a.index++
		return token, nil
	}
	// we are at the key, so return the key
	token := Token{
		Typ: AttrKeyToken,
		Val: a.attrs[a.index].key,
	}
	// move onto the value of this key
	a.atValue = true
	return token, nil
}

type attr struct {
	key string
	val string
}

func newAttr(a xml.Attr) *attr {
	return &attr{
		key: a.Name.Local,
		val: a.Value,
	}
}
