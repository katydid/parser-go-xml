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

package token

import (
	"testing"

	"katydid.org.za/go/parser-go-xml/xml/scan"
)

func TestPersonManual(t *testing.T) {
	personStr := `<Person>
	<Name fullname=false>Norris</Name>
	<Addresses>
		<Number>456</Number>
		<Street>TheStreet</Street>
	</Addresses>
	<Weight>102.3</Weight>
	<XXX_unrecognized/>
</Person>`
	x := NewTokenizer(WithBuffer([]byte(personStr)))
	expect(t, x.Next, scan.StartKind)
	expectStr(t, x, "Person")
	expect(t, x.Next, scan.CharKind)
	expectStr(t, x, "\n\t")
	expect(t, x.Next, scan.StartKind)
	expectStr(t, x, "Name")
	expect(t, x.Next, scan.AttrKeyKind)
	expectStr(t, x, "fullname")
	expect(t, x.Next, scan.AttrValueKind)
	expectFalse(t, x)
	expect(t, x.Next, scan.CharKind)
	expectStr(t, x, "Norris")
	expect(t, x.Next, scan.EndKind)
	expectStr(t, x, "Name")
	expect(t, x.Next, scan.CharKind)
	expectStr(t, x, "\n\t")
	expect(t, x.Next, scan.StartKind)
	expectStr(t, x, "Addresses")
	expect(t, x.Next, scan.CharKind)
	expectStr(t, x, "\n\t\t")
	expect(t, x.Next, scan.StartKind)
	expectStr(t, x, "Number")
	expect(t, x.Next, scan.CharKind)
	expectInt(t, x, 456)
	expect(t, x.Next, scan.EndKind)
	expectStr(t, x, "Number")
	expect(t, x.Next, scan.CharKind)
	expectStr(t, x, "\n\t\t")
	expect(t, x.Next, scan.StartKind)
	expectStr(t, x, "Street")
	expect(t, x.Next, scan.CharKind)
	expectStr(t, x, "TheStreet")
	expect(t, x.Next, scan.EndKind)
	expectStr(t, x, "Street")
	expect(t, x.Next, scan.CharKind)
	expectStr(t, x, "\n\t")
	expect(t, x.Next, scan.EndKind)
	expectStr(t, x, "Addresses")
	expect(t, x.Next, scan.CharKind)
	expectStr(t, x, "\n\t")
	expect(t, x.Next, scan.StartKind)
	expectStr(t, x, "Weight")
	expect(t, x.Next, scan.CharKind)
	expectFloat(t, x, 102.3)
	expect(t, x.Next, scan.EndKind)
	expectStr(t, x, "Weight")
	expect(t, x.Next, scan.CharKind)
	expectStr(t, x, "\n\t")
	expect(t, x.Next, scan.StartKind)
	expectStr(t, x, "XXX_unrecognized")
	expect(t, x.Next, scan.EndKind)
	expectStr(t, x, "XXX_unrecognized")
	expect(t, x.Next, scan.CharKind)
	expectStr(t, x, "\n")
	expect(t, x.Next, scan.EndKind)
	expectStr(t, x, "Person")
	expectEOF(t, x.Next)
}
