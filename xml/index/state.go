//  Copyright 2025 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package index

import "github.com/katydid/parser-go/parse"

type state struct {
	kind  stateKind
	index int64
	hint  parse.Hint
}

type stateKind byte

const startState = stateKind(0)

const arrayOpenState = stateKind('[')
const arrayIndexState = stateKind('i')
const arrayElemState = stateKind('e')

const objectOpenState = stateKind('{')
const objectKeyState = stateKind('k')
const objectValueState = stateKind('v')

const endState = stateKind('$')
