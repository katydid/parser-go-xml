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

type options struct {
	attrPrefix string
	elemPrefix string
	textPrefix string
}

func newOptions(opts ...Option) *options {
	x := &options{}
	for _, option := range opts {
		option(x)
	}
	return x
}

// Option is used set options when creating a new XMLParser
type Option func(x *options)

// WithAttrPrefix specifies the prefix which will be added to attributes returned by the parser.
func WithAttrPrefix(a string) func(x *options) {
	return func(x *options) {
		x.attrPrefix = a
	}
}

// WithElemPrefix specifies the prefix which will be added to elements returned by the parser.
func WithElemPrefix(e string) func(x *options) {
	return func(x *options) {
		x.elemPrefix = e
	}
}

// WithTextPrefix specifies the prefix which will be added to text returned by the parser.
func WithTextPrefix(e string) func(x *options) {
	return func(x *options) {
		x.textPrefix = e
	}
}
