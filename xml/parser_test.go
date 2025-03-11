//  Copyright 2015 Walter Schulze
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

package xml

import (
	"encoding/json"
	"testing"

	"github.com/katydid/parser-go/parser/debug"
)

func testXML(t *testing.T, s string) {
	x := newParser(t, s)
	m, err := debug.Parse(debug.NewLogger(x, debug.NewLineLogger()))
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

func newParser(t *testing.T, s string) XMLParser {
	t.Helper()
	x := NewXMLParser()
	if err := x.Init([]byte(s)); err != nil {
		t.Fatal(err)
	}
	return x
}

func TestExample(t *testing.T) {
	example := `
		<Top>
			<Name>Katydid</Name>
			<Dragons alive="false">
				<Fire>true</Fire>
			</Dragons>
			<Empty></Empty>
		</Top>`
	testXML(t, example)
}

func TestPudding(t *testing.T) {
	pudding := `
	<FinanceJudo>
	<SaladWorry><MeasureGrade></MeasureGrade><MagazineFrame>a</MagazineFrame><MagazineFrame>b</MagazineFrame>
		<XrayPilot><AnkleCoat>2</AnkleCoat><XXX_unrecognized></XXX_unrecognized></XrayPilot><XXX_unrecognized></XXX_unrecognized>
	</SaladWorry>
	<RumourSpirit>1</RumourSpirit><XXX_unrecognized></XXX_unrecognized>
	</FinanceJudo>
	`
	testXML(t, pudding)
}

func TestAB(t *testing.T) {
	testXML(t, `<A>B</A>`)
}

func TestPersonWalk(t *testing.T) {
	personStr := `<Person>
		<Name>Robert</Name>
		<Addresses>
				<Number>456</Number>
				<Street>TheStreet</Street>
		</Addresses>
		<Telephone>0127897897</Telephone>
		<XXX_unrecognized/>
	</Person>`
	testXML(t, personStr)
	// [{"Label":"Person","Children":[{"Label":"\n\t\t","Children":null},{"Label":"Name","Children":[{"Label":"Robert","Children":null}]},{"Label":"\n\t\t","Children":null},{"Label":"Addresses","Children":[{"Label":"\n\t\t\t\t","Children":null},{"Label":"Number","Children":[{"Label":"456","Children":null}]},{"Label":"\n\t\t\t\t","Children":null},{"Label":"Street","Children":[{"Label":"TheStreet","Children":null}]},{"Label":"\n\t\t","Children":null}]},{"Label":"\n\t\t","Children":null},{"Label":"Telephone","Children":[{"Label":"127897897","Children":null}]},{"Label":"\n\t\t","Children":null},{"Label":"XXX_unrecognized","Children":[]},{"Label":"\n\t","Children":null}]}]
}

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
	x := newParser(t, personStr)
	expectNoErr(t, x.Next)
	expect(t, x.String, "Person")
	x.Down()
	expectNoErr(t, x.Next)
	expect(t, x.String, "\n\t")
	expectNoErr(t, x.Next)
	expect(t, x.String, "Name")
	x.Down()
	expectNoErr(t, x.Next)
	expect(t, x.String, "Robert")
	expectEOF(t, x.Next)
	x.Up()
	expectNoErr(t, x.Next)
	expect(t, x.String, "\n\t")
	expectNoErr(t, x.Next)
	expect(t, x.String, "Addresses")
	expectNoErr(t, x.Next)
	expect(t, x.String, "\n\t")
	expectNoErr(t, x.Next)
	expect(t, x.String, "Telephone")
	x.Down()
	expectNoErr(t, x.Next)
	expect(t, x.String, "0127897897")
	expectEOF(t, x.Next)
	x.Up()
	expectNoErr(t, x.Next)
	expect(t, x.String, "\n\t")
	expectNoErr(t, x.Next)
	expect(t, x.String, "XXX_unrecognized")
	x.Down()
	expectEOF(t, x.Next)
	x.Up()
	expectNoErr(t, x.Next)
	expect(t, x.String, "\n")
	expectEOF(t, x.Next)
}
