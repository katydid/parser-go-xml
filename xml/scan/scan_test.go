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
	x := NewScanner(WithBuffer([]byte(personStr)))
	expect(t, x.Next, StartKind, "Person")
	expect(t, x.Next, CharKind, "\n\t")
	expect(t, x.Next, StartKind, "Name")
	expect(t, x.Next, CharKind, "Robert")
	expect(t, x.Next, EndKind, "Name")
	expect(t, x.Next, CharKind, "\n\t")
	expect(t, x.Next, StartKind, "Addresses")
	expect(t, x.Next, CharKind, "\n\t\t")
	expect(t, x.Next, StartKind, "Number")
	expect(t, x.Next, CharKind, "456")
	expect(t, x.Next, EndKind, "Number")
	expect(t, x.Next, CharKind, "\n\t\t")
	expect(t, x.Next, StartKind, "Street")
	expect(t, x.Next, CharKind, "TheStreet")
	expect(t, x.Next, EndKind, "Street")
	expect(t, x.Next, CharKind, "\n\t")
	expect(t, x.Next, EndKind, "Addresses")
	expect(t, x.Next, CharKind, "\n\t")
	expect(t, x.Next, StartKind, "Telephone")
	expect(t, x.Next, CharKind, "0127897897")
	expect(t, x.Next, EndKind, "Telephone")
	expect(t, x.Next, CharKind, "\n\t")
	expect(t, x.Next, StartKind, "XXX_unrecognized")
	expect(t, x.Next, EndKind, "XXX_unrecognized")
	expect(t, x.Next, CharKind, "\n")
	expect(t, x.Next, EndKind, "Person")
	expectEOF(t, x.Next)
}

func TestAttrManual(t *testing.T) {
	personStr := `<Person name="Robert"><Address number=456 street="TheStreet"/></Person>`
	x := NewScanner(WithBuffer([]byte(personStr)))
	expect(t, x.Next, StartKind, "Person")
	expect(t, x.Next, AttrKeyKind, "name")
	expect(t, x.Next, AttrValueKind, "Robert")
	expect(t, x.Next, StartKind, "Address")
	expect(t, x.Next, AttrKeyKind, "number")
	expect(t, x.Next, AttrValueKind, "456")
	expect(t, x.Next, AttrKeyKind, "street")
	expect(t, x.Next, AttrValueKind, "TheStreet")
	expect(t, x.Next, EndKind, "Address")
	expect(t, x.Next, EndKind, "Person")
	expectEOF(t, x.Next)
}
