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

type state byte

const startState = state(0)

// The underlying parser has the following state transitions for elements:

// startedState/arrayOpenedState -> scan.StartKind
// p.down(objectOpenedState)
// return parse.ObjectOpenHint, nil

// objectOpenedState
// p.state = objectKeyedState
// return parse.KeyHint, nil

// objectKeyedState
// p.state = objectValuedState
// p.down(arrayOpenedState)
// return parse.ArrayOpenHint, nil

// arrayOpenedState -> scan.EndKind
// p.up()
// return parse.ArrayCloseHint, nil

// objectValuedState
// p.up()
// return parse.ObjectCloseHint, nil

// This implies that
// we need to go to elemTagOpenState only when
// ScanKind = scan.StartKind AND Next's Hint = parse.ObjectOpenHint
// we need to go to elemTagCloseState only whe
// ScanKind = scan.EndKind AND Next's Hint = parse.ObjectCloseHint

const elemTagOpenState = state('<')

const elemTagKeyState = state('e')

const elemTagCloseState = state('>')

// The underlying parser has the following state transitions for attributes:

// arrayOpenedState -> scan.AttrKeyKind
// p.down(attrOpenedState)
// return parse.ObjectOpenHint, nil

// attrOpenedState
// p.state = attrKeyedState
// return parse.KeyHint, nil

// attrKeyedState -> scan.AttrValueKind
// p.state = attrValuedState
// return parse.ValueHint, nil

// attrValuedState
// p.up()
// return parse.ObjectCloseHint, nil

// This implies that
// we need to go to attrTagOpenState only when
// ScanKind = scan.AttrKeyKind AND Next's Hint = parse.ObjectOpenHint
// we need to go to attrTagCloseState only whe
// ScanKind = scan.AttrValueKind AND Next's Hint = parse.ObjectCloseHint

const attrTagOpenState = state('{')

const attrTagKeyState = state('a')

const attrTagCloseState = state('}')

// The underlying parser has the following state transitions for text elements:

// startedState/arrayOpenedState -> scan.CharKind
// return parse.ValueHint, nil

// This implies for text elements that
// we need to go to textTagOpenState only when
// ScanKind = scan.CharKind AND Next's Hint = parse.ValueHint
// and that the rest of the states do not need a stack.

// The underlying parser has the following state transitions for attribute text values:

// attrKeyedState -> scan.AttrValueKind
// p.state = attrValuedState
// return parse.ValueHint, nil

// This implies for attribute text values that
// we need to go to textTagOpenState only when
// ScanKind = scan.AttrValueKind AND Next's Hint = parse.ValueHint
// and that the rest of the states do not need a stack.

const textTagOpenState = state('(')

const textTagKeyState = state('t')

const textTagCloseState = state(')')

const endState = state('$')
