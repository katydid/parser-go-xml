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

type Kind byte

const UnknownKind = Kind(0)
const StartKind = Kind('<')
const AttrKeyKind = Kind('k')
const AttrValueKind = Kind('v')
const CharKind = Kind('c')
const EndKind = Kind('>')

func (k Kind) String() string {
	switch k {
	case UnknownKind:
		return "UnknownKind"
	case StartKind:
		return "StartKind"
	case AttrKeyKind:
		return "AttrKeyKind"
	case AttrValueKind:
		return "AttrValueKind"
	case CharKind:
		return "CharKind"
	case EndKind:
		return "EndKind"
	}
	panic("unreachable")
}
