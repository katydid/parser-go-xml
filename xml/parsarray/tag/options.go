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

type options struct {
	tagAttrs bool
	tagElems bool
	tagTexts bool

	// deprecated
	attrPrefix []byte
	// deprecated
	elemPrefix []byte
	// deprecated
	textPrefix []byte
}

func newOptions(opts ...Option) *options {
	x := &options{}
	for _, option := range opts {
		option(x)
	}
	return x
}

type Option func(x *options)

// Deprecated: use WithAttrTab
// WithAttrPrefix specifies the prefix which will be added to attributes returned by the parser.
func WithAttrPrefix(a string) func(x *options) {
	return func(x *options) {
		x.attrPrefix = []byte(a)
	}
}

// WithAttrTag tags each attribute key with an "attr" key.
// For example:
// Without the option: <a k="v"/> will parse as [{"a": [{"k": "v"}]}]
// With the option: <a k="v"/> will parse as [{"a": [{"attr": {"k": "v"}}]}]
func WithAttrTag() func(x *options) {
	return func(x *options) {
		x.tagAttrs = true
	}
}

// Deprecated: use WithElemTag
// WithElemPrefix specifies the prefix which will be added to elements returned by the parser.
func WithElemPrefix(e string) func(x *options) {
	return func(x *options) {
		x.elemPrefix = []byte(e)
	}
}

// WithElemTag tags each element with an "elem" key.
// For example:
// Without the option: <a/> will parse as [{"a": []}]
// With the option: <a/> will parse as [{"elem": {"a": []}}]
func WithElemTag() func(x *options) {
	return func(x *options) {
		x.tagElems = true
	}
}

// Deprecated: use WithTextTag
// WithTextPrefix specifies the prefix which will be added to text returned by the parser.
func WithTextPrefix(e string) func(x *options) {
	return func(x *options) {
		x.textPrefix = []byte(e)
	}
}

// WithTextTag tags each text element and attribute value with a "text" key.
// For example:
// Without the option: <a>t</a> will parse as [{"a": ["t"]}]
// With the option: <a>t</a> will parse as [{"a": [{"text": "t"}]}]
// Without the option: <a k="v"/> will parse as [{"a": [{"k": "v"}]}]
// With the option: <a k="v"/> will parse as [{"a": [{"k": {"text": "v"}}]}]
func WithTextTag() func(x *options) {
	return func(x *options) {
		x.tagTexts = true
	}
}
