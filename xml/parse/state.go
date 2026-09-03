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

package parse

type state byte

const startState = state(0)

const inElemState = state('{')

const atFieldState = state('F')

const isLeafState = state('l')

const atAttributeKeyState = state('k')

const atAttributeValueState = state('v')

const endState = state('$')

func (s state) String() string {
	switch s {
	case startState:
		return "startState"
	case inElemState:
		return "inElemState"
	case atFieldState:
		return "atFieldState"
	case isLeafState:
		return "isLeafState"
	case atAttributeKeyState:
		return "atAttributeKeyState"
	case atAttributeValueState:
		return "atAttributeValueState"
	case endState:
		return "endState"
	}
	panic("unreachable")
}
