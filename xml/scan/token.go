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

type TokenType byte

const UnknownToken = TokenType(0)
const StartToken = TokenType(1)
const AttrKeyToken = TokenType(2)
const AttrValueToken = TokenType(3)
const CharToken = TokenType(4)
const EndToken = TokenType(5)

type Token struct {
	Typ TokenType
	Val string
}
