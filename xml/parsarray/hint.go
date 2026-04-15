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

package parse

// Hint of the token that is parsed.
// This is represented by one for following bytes: {k}[v]
type Hint byte

const UnknownHint = Hint(0)

func (h Hint) IsUnknown() bool {
	return h == UnknownHint
}

const ObjectOpenHint = Hint('{')

func (h Hint) IsObjectOpen() bool {
	return h == ObjectOpenHint
}

const KeyHint = Hint('k')

func (h Hint) IsKey() bool {
	return h == KeyHint
}

const ValueHint = Hint('v')

func (h Hint) IsValue() bool {
	return h == ValueHint
}

const ObjectCloseHint = Hint('}')

func (h Hint) IsObjectClose() bool {
	return h == ObjectCloseHint
}

const ArrayOpenHint = Hint('[')

func (h Hint) IsArrayOpen() bool {
	return h == ArrayOpenHint
}

const ArrayCloseHint = Hint(']')

func (h Hint) IsArrayClose() bool {
	return h == ArrayCloseHint
}

func (h Hint) String() string {
	switch h {
	case UnknownHint:
		return "unknown"
	case ArrayOpenHint:
		return "arrayOpen"
	case ValueHint:
		return "value"
	case ArrayCloseHint:
		return "arrayClose"
	case ObjectOpenHint:
		return "objectOpen"
	case KeyHint:
		return "key"
	case ObjectCloseHint:
		return "objectClose"
	}
	panic("unreachable")
}
