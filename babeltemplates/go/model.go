{{template "SIMPLECOMMENTS" .Comments }}
package {{package}}

// *** AUTO-GENERATED FILE - DO NOT MODIFY ***
// *** Generated from {{.Filename}} ***

import ({{if modelUsesType "decimal"}}
	"math/big"{{end}}{{if modelUsesType "datetime"}}
	"time"{{end}}
{{range imports}}	"{{.}}"
{{end}})

{{range .Enums}}{{$nm := .Name}}{{setindent ""}}{{template "COMMENTS" .Comments }}type {{$nm}} string

// Values for type {{$nm}}
const (
{{range $i,$v := .Values}}	{{$nm}}{{.Name}} = "{{.Name}}"
{{end}})

// Get{{$nm}} returns the {{$nm}} for a given integer value.
func Get{{$nm}}(value int) ({{$nm}}, bool) {
	switch value {
{{range $i,$v := .Values}}	case {{formatValue .}}:
		return {{$nm}}{{.Name}}, true
{{end}}	default:
		return "", false
	}
}

// Value converts a string to the enumerated type's integer value.
// Returns true when the string is found and false when not.
func (s *{{$nm}}) Value() (int, bool) {
	if s == nil {
		return 0, false
	}
	switch *s {
{{range $i,$v := .Values}}	case {{$nm}}{{.Name}}:
		return {{formatValue .}}, true
{{end}}	default:
		return 0, false
	}
}

{{end}}
{{range .Consts}}{{$nm := .Name}}{{setindent ""}}{{template "COMMENTS" .Comments }}const (
{{range .Values}}	{{$nm}}{{.Name}} = {{formatValue .}}
{{end}})

{{end}}
{{range $is, $xs :=  .Structs}}
{{setindent ""}}{{template "COMMENTS" .Comments }}type {{.Name}} struct {{"{"}}{{if .Extends}}
	{{.Extends}}
{{end}}{{range .Fields}}
{{setindent "\t"}}{{template "COMMENTS" .Comments }}	{{toPascalCase .Name}} {{formatType .Type}} `json:"{{.Name}}{{serializerOptions .Type}}"`
{{end}}
}

// Init sets default values for a {{.Name}}
func (obj *{{.Name}}) Init() *{{.Name}} {{"{"}}{{range .Fields}}{{if .Initializer}}
	obj.{{toPascalCase .Name}} = new({{notptr (formatType .Type)}})
	*obj.{{toPascalCase .Name}} = {{formatValue .Initializer}}{{end}}{{if .Type.IsList}}
	obj.{{toPascalCase .Name}} = make({{formatType .Type}}, 0){{end}}{{if .Type.IsMap}}
	obj.{{toPascalCase .Name}} = make({{formatType .Type}}, 0){{end}}{{end}}
	return obj
}
{{end}}