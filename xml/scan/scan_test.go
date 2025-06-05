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

import "testing"

func TestPersonManual(t *testing.T) {
	personStr := `<Person>
	<Name>Robert</Name>
	<Addresses>
		<Number>456</Number>
		<Street>TheStreet</Street>
	</Addresses>
	<Telephone>0127897897</Telephone>
	<XXX_unrecognized/>
</Person>`
	x := NewScanner([]byte(personStr))
	expect(t, x.Next, Token{StartToken, "Person"})
	expect(t, x.Next, Token{CharToken, "\n\t"})
	expect(t, x.Next, Token{StartToken, "Name"})
	expect(t, x.Next, Token{CharToken, "Robert"})
	expect(t, x.Next, Token{EndToken, "Name"})
	expect(t, x.Next, Token{CharToken, "\n\t"})
	expect(t, x.Next, Token{StartToken, "Addresses"})
	expect(t, x.Next, Token{CharToken, "\n\t\t"})
	expect(t, x.Next, Token{StartToken, "Number"})
	expect(t, x.Next, Token{CharToken, "456"})
	expect(t, x.Next, Token{EndToken, "Number"})
	expect(t, x.Next, Token{CharToken, "\n\t\t"})
	expect(t, x.Next, Token{StartToken, "Street"})
	expect(t, x.Next, Token{CharToken, "TheStreet"})
	expect(t, x.Next, Token{EndToken, "Street"})
	expect(t, x.Next, Token{CharToken, "\n\t"})
	expect(t, x.Next, Token{EndToken, "Addresses"})
	expect(t, x.Next, Token{CharToken, "\n\t"})
	expect(t, x.Next, Token{StartToken, "Telephone"})
	expect(t, x.Next, Token{CharToken, "0127897897"})
	expect(t, x.Next, Token{EndToken, "Telephone"})
	expect(t, x.Next, Token{CharToken, "\n\t"})
	expect(t, x.Next, Token{StartToken, "XXX_unrecognized"})
	expect(t, x.Next, Token{EndToken, "XXX_unrecognized"})
	expect(t, x.Next, Token{CharToken, "\n"})
	expect(t, x.Next, Token{EndToken, "Person"})
}

func TestAttrManual(t *testing.T) {
	personStr := `<Person name="Robert"><Address number=456 street="TheStreet"/></Person>`
	x := NewScanner([]byte(personStr))
	expect(t, x.Next, Token{StartToken, "Person"})
	expect(t, x.Next, Token{AttrKeyToken, "name"})
	expect(t, x.Next, Token{AttrValueToken, "Robert"})
	expect(t, x.Next, Token{StartToken, "Address"})
	expect(t, x.Next, Token{AttrKeyToken, "number"})
	expect(t, x.Next, Token{AttrValueToken, "456"})
	expect(t, x.Next, Token{AttrKeyToken, "street"})
	expect(t, x.Next, Token{AttrValueToken, "TheStreet"})
	expect(t, x.Next, Token{EndToken, "Address"})
	expect(t, x.Next, Token{EndToken, "Person"})
}
