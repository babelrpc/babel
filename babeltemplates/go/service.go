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
{{end}}
{{setindent ""}}{{template "COMMENTS" .Comments }}
type I{{.Name}} interface { {{range .Methods}}
{{setindent "\t"}}{{template "METHODCOMMENTS" .}}{{indent}}{{toPascalCase .Name}}({{range $i, $x := .Parameters}}{{if $i}}, {{end}}{{.Name}} {{formatType .Type}}{{end}}) {{if formatType .Returns}}({{formatType .Returns}}, error){{else}}error{{end}} 
{{end}}}{{end}}
