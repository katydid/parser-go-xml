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

import "testing"

func TestPersonPeekSkipSpace(t *testing.T) {
	personStr := `<Person>
	<Name>Robert</Name>
	<Addresses>
		<Number>456</Number>
		<Street>TheStreet</Street>
	</Addresses>
	<Telephone>0127897897</Telephone>
	<XXX_unrecognized/>
</Person>`
	x := NewScanner(WithBuffer([]byte(personStr)), WithSkipSpace())
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "Person")
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "Name")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "Robert")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "Name")
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "Addresses")
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "Number")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "456")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "Number")
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "Street")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "TheStreet")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "Street")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "Addresses")
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "Telephone")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "0127897897")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "Telephone")
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "XXX_unrecognized")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "XXX_unrecognized")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "Person")
	expectPeekEOF(t, x.Peek)
	expectEOF(t, x.Next)
}

func TestPersonPeekManual(t *testing.T) {
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
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "Person")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "\n\t")
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "Name")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "Robert")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "Name")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "\n\t")
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "Addresses")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "\n\t\t")
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "Number")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "456")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "Number")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "\n\t\t")
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "Street")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "TheStreet")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "Street")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "\n\t")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "Addresses")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "\n\t")
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "Telephone")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "0127897897")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "Telephone")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "\n\t")
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "XXX_unrecognized")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "XXX_unrecognized")
	expectPeek(t, x.Peek, CharKind)
	expect(t, x.Next, CharKind, "\n")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "Person")
	expectPeekEOF(t, x.Peek)
	expectEOF(t, x.Next)
}

func TestAttrPeekManual(t *testing.T) {
	personStr := `<Person name="Robert"><Address number=456 street="TheStreet"/></Person>`
	x := NewScanner(WithBuffer([]byte(personStr)))
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "Person")
	expectPeek(t, x.Peek, AttrKeyKind)
	expect(t, x.Next, AttrKeyKind, "name")
	expectPeek(t, x.Peek, AttrValueKind)
	expect(t, x.Next, AttrValueKind, "Robert")
	expectPeek(t, x.Peek, StartKind)
	expect(t, x.Next, StartKind, "Address")
	expectPeek(t, x.Peek, AttrKeyKind)
	expect(t, x.Next, AttrKeyKind, "number")
	expectPeek(t, x.Peek, AttrValueKind)
	expect(t, x.Next, AttrValueKind, "456")
	expectPeek(t, x.Peek, AttrKeyKind)
	expect(t, x.Next, AttrKeyKind, "street")
	expectPeek(t, x.Peek, AttrValueKind)
	expect(t, x.Next, AttrValueKind, "TheStreet")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "Address")
	expectPeek(t, x.Peek, EndKind)
	expect(t, x.Next, EndKind, "Person")
	expectPeekEOF(t, x.Peek)
	expectEOF(t, x.Next)
}
