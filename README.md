## parser-go-xml

Parser for XML in Go

```go
import (
    "katydid.org.za/go/parser-go/cast"
	"katydid.org.za/go/parser-go/cp"
    "katydid.org.za/go/parser-go-xml/xml"
    "katydid.org.za/go/parser-go-xml/xml/parse"
)

func main() {
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
	println(val) // Katydid
}

func getName(xmlstr string) (string, error) {
	p := xml.NewParser()
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
```