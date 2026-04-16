// Copyright 2026 Walter Schulze
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

type options struct {
	buf       []byte
	skipSpace bool
}

func newOptions(opts ...Option) *options {
	o := &options{
		skipSpace: false,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Option is used set options when creating a new JSON Parser.
type Option func(*options)

// WithBuffer passes in a buffer to parse.
func WithBuffer(buf []byte) func(*options) {
	return func(o *options) {
		o.buf = buf
	}
}

// WithSkipSpace skips all Value nodes that only contain space characters.
func WithSkipSpace() func(*options) {
	return func(o *options) {
		o.skipSpace = true
	}
}
