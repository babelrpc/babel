{{define "COMMENTS"}}
{{$cmts := expandComments .}}{{if len $cmts}}
{{indent}}/**
{{range .}}{{indent}} * {{.}}
{{end}}{{indent}} */{{end}}{{end}}

{{define "SIMPLECOMMENTS"}}{{range .}}{{indent}}//{{.}}
{{end}}{{end}}

{{define "METHODCOMMENTS"}}{{$cmts := expandComments .Comments}}
{{indent}}/**
{{range $cmts}}{{indent}} * {{.}}{{end}}
{{range .Parameters}}{{indent}} * @param { {{.Name}} } {{range $i,$m := .Comments}}{{if $i}}
{{indent}}{{end}}{{.}}{{end}}
{{end}}{{indent}} */
{{end}}

{{define "ATTRS"}}{{$attrs := filterAttrs .}}{{if len $attrs}}{{indent}}[{{range $i, $x := $attrs}}{{if $i}}, {{end}}{{.Name}}{{if len .Parameters}}({{range $j, $y := .Parameters}}{{if $j}}, {{end}}{{if .Name}}{{.Name}} = {{end}}{{formatValue .}}{{end}}){{end}}{{end}}]
{{end}}{{end}}