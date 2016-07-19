package gen

import (
	"text/template"

	"github.com/xeger/mongoose/parse"
)

const gomutiItem = `
{{$typename := .Interface.Name | printf "Mock%s" }}type {{$typename}} struct {
	Stub bool
	Mock gomuti.Mock
}

{{$locl := .Package.Name}}{{$res := .Resolver}}{{range .Interface.Methods}}
func (m {{$typename}}) {{.Name}}{{.Params.Tuple $locl $res}}{{$rtuple := .Results.Tuple $locl $res}}{{if gt .Results.Len 0}} {{$rtuple}}{{end}} {
	{{$pnames := .Params.NameList}}{{$ptypes := (.Params.TypeList $locl $res)}}ret := m.Mock.Delegate("{{.Name}}",{{.Params.NameList}})
	if ret == nil {
		if m.Stub {
			return{{if gt .Results.Len 0}} {{.Results.ZeroList $locl $res}}{{end}}
		}
		panic("{{$typename}}: unexpected call to {{.Name}}")
	}
	{{range $idx, $typ := .Results}}
	var r{{$idx}} {{$typ.ShortName $locl $res}}
	if ret[{{$idx}}] == nil {
		r{{$idx}} = {{$typ.ZeroValue $locl $res}}
	} else {
		r{{$idx}} = ret[{{$idx}}].({{$typ.ShortName $locl $res}})
	}
	{{end}}
	return {{.Results.NameList}}
}
{{end}}
`

// NewGomutiRenderer creates a code generator using github.com/xeger/gomuti.
// The mock type contains a gomuti.Mock and can be programmed using the
// gomuti.Allow() method.
func NewGomutiRenderer() Renderer {
	r := parse.NewResolver()
	r.Import("gomuti", "github.com/xeger/gomuti")
	tr := &templateRenderer{}
	tr.Resolver = r
	tr.Header = template.New("gomutiHeader")
	tr.Header.Parse(templateHeader)
	tr.Item = template.New("gomutiItem")
	tr.Item.Parse(gomutiItem)
	return tr
}
