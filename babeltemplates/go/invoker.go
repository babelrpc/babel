{{template "SIMPLECOMMENTS" .Comments }}
package {{package}}

// *** AUTO-GENERATED FILE - DO NOT MODIFY ***
// *** Generated from {{.Filename}} ***

import ({{if serviceUsesType "decimal"}}
	"math/big"{{end}}{{if serviceUsesType "datetime"}}
	"time"{{end}}
{{range imports}}	"{{.}}"
{{end}})

{{range $k, $s := .Services}}{{if $k}}
{{end}}{{range .Methods}}
{{setindent ""}}// {{$s.Name}}{{.Name}}Request is the request structure used for invoking the {{.Name}} method on the {{$s.Name}} service.
type {{$s.Name}}{{.Name}}Request struct {{"{"}}
{{range $i, $x := .Parameters}}
{{setindent "\t"}}{{template "COMMENTS" .Comments }}	{{toPascalCase .Name}} {{formatType .Type}} `json:"{{.Name}}{{serializerOptions .Type}}"`
{{end}}
}

// Init sets default values for a {{.Name}}Request
func (obj *{{$s.Name}}{{.Name}}Request) Init() *{{$s.Name}}{{.Name}}Request {{"{"}}{{range .Parameters}}{{if .Initializer}}
	obj.{{toPascalCase .Name}} = new({{notptr (formatType .Type)}})
	*obj.{{toPascalCase .Name}} = {{formatValue .Initializer}}{{end}}{{if .Type.IsList}}
	obj.{{toPascalCase .Name}} = make({{formatType .Type}}, 0){{end}}{{if .Type.IsMap}}
	obj.{{toPascalCase .Name}} = make({{formatType .Type}}, 0){{end}}{{end}}
	return obj
}

{{setindent ""}}// {{$s.Name}}{{.Name}}Response is the response structure used for invoking the {{.Name}} method on the {{$s.Name}} service.
type {{$s.Name}}{{.Name}}Response struct {{"{"}}
{{if formatType .Returns}}
{{setindent "\t"}}	Value {{formatType .Returns}} `json:"Value{{serializerOptions .Returns}}"`
{{end}}
}

// Init sets default values for a {{.Name}}Response
func (obj *{{$s.Name}}{{.Name}}Response) Init() *{{$s.Name}}{{.Name}}Response {{"{"}}{{if formatType .Returns}}{{if .Returns.IsList}}
	obj.Value = make({{formatType .Returns}}, 0){{end}}{{if .Returns.IsMap}}
	obj.Value = make({{formatType .Returns}}, 0){{end}}{{end}}
	return obj
}

{{end}}
{{end}}

{{range $k, $s := .Services}}{{if $k}}
{{end}}{{setindent ""}}{{template "COMMENTS" .Comments }}type {{$s.Name}} struct {
	SvcObj I{{$s.Name}} `json:"-"`
}
{{range .Methods}}
{{setindent ""}}{{template "COMMENTS" .Comments }}func (s *{{$s.Name}}) {{.Name}}(req *{{$s.Name}}{{.Name}}Request, rsp *{{$s.Name}}{{.Name}}Response) error {
	{{if formatType .Returns}}response, {{end}}err := s.SvcObj.{{.Name}}({{range $pk, $ps := .Parameters}}{{if $pk}}, {{end}}req.{{toPascalCase $ps.Name}}{{end}})
{{if formatType .Returns}}	if err == nil {
		rsp.Value = response
	}{{end}}
	return err
}
{{end}}
{{end}}

