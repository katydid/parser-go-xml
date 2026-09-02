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

	xmlparse "katydid.org.za/go/parser-go-xml/xml/parse"
	"katydid.org.za/go/parser-go/cast"
	"katydid.org.za/go/parser-go/cp"
	"katydid.org.za/go/parser-go/expect"
	"katydid.org.za/go/parser-go/hedge"
	"katydid.org.za/go/parser-go/parse"
	"katydid.org.za/go/parser-go/parse/log"
)

func testXML(t *testing.T, s string) {
	x := NewParser()
	x.Init([]byte(s))
	m, err := hedge.ParseInto(log.WrapParserWithInit(x))
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

func TestExample(t *testing.T) {
	xmlstr := `
		<Top>
			<Name>Katydid</Name>
			<Dragons alive="false">
				<Fire>true</Fire>
			</Dragons>
		</Top>`
	val, err := getName(xmlstr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}

func getName(xmlstr string) (string, error) {
	p := NewParser()
	p.Init([]byte(xmlstr))

	p.Next()              // parse.EnterHint
	p.Next()              // parse.FieldHint Token: StringKind, "Top"
	p.Next()              // parse.EnterHint
	hint, err := p.Next() // parse.FieldHint Token: StringKind, "Name"
	if err != nil {
		panic(err)
	}
	if hint != parse.FieldHint {
		panic("expected field hint")
	}
	kind, val, err := p.Token()
	if err != nil {
		panic(err)
	}
	if kind != parse.StringKind {
		panic("expected string kind")
	}
	var s string
	cast.ToStringPtr(val, &s)
	if s != "Name" {
		panic("expected Name")
	}
	hint, err = p.Next() // parse.ValueHint Token: StringKind, "Katydid"
	if err != nil {
		panic(err)
	}
	if hint != parse.ValueHint {
		panic("expected field hint")
	}
	kind, val, err = p.Token()
	if err != nil {
		panic(err)
	}
	if kind != parse.StringKind {
		panic("expected string kind")
	}
	return cp.ToString(val), nil
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

func newParser(s string) xmlparse.Parser {
	x := NewRAWParser()
	x.Init([]byte(s))
	return x
}

func TestPersonManualSkipAddresses(t *testing.T) {
	personStr := `<Person>
	<Name>Robert</Name>
	<Addresses>
		<Number>456</Number>
		<Street>TheStreet</Street>
	</Addresses>
	<Telephone>0127897897</Telephone>
	<XXX_unrecognized/>
</Person>`
	p := newParser(personStr)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "Person")
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "\n\t")

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "Name")
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "Robert")

	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "\n\t")

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "Addresses")
	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "\n\t")

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "Telephone")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 127897897)

	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "\n\t")

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "XXX_unrecognized")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.LeaveHint)

	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "\n")

	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestPersonManualDownAddresses(t *testing.T) {
	personStr := `<Person>
	<Name>Robert</Name>
	<Addresses>
		<Number>456</Number>
		<Street>TheStreet</Street>
	</Addresses>
	<Telephone>0127897897</Telephone>
	<XXX_unrecognized/>
</Person>`
	p := newParser(personStr)
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "Person")
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "\n\t")

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "Name")
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "Robert")

	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "\n\t")

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "Addresses")
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "\n\t\t")

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "Number")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 456)

	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "\n\t\t")

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "Street")
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "TheStreet")

	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "\n\t")

	expect.Hint(t, p, parse.LeaveHint)

	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "\n\t")

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "Telephone")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 127897897)

	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "\n\t")

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "XXX_unrecognized")
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.LeaveHint)

	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "\n")

	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}

func TestAttrManual(t *testing.T) {
	personStr := `<Person name="Robert"><Address number=456 street="TheStreet"/></Person>`
	p := newParser(personStr)
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "Person")
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "name")
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "Robert")

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "Address")
	expect.Hint(t, p, parse.EnterHint)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "number")
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 456)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "street")
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "TheStreet")

	expect.Hint(t, p, parse.LeaveHint)

	expect.Hint(t, p, parse.LeaveHint)
	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}
