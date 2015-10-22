{{define "COMMENTS"}}
{{$cmts := expandComments .}}{{if len $cmts}}
{{indent}}/**
{{range $cmts}}{{indent}} * {{.}}
{{end}}{{indent}} */{{end}}{{end}}
{{define "SIMPLECOMMENTS"}}{{range .}}{{indent}}//{{.}}
{{end}}{{end}}
{{define "ATTRS"}}{{$attrs := filterAttrs .}}{{if len $attrs}}{{indent}}{{range $i, $x := $attrs}}{{if $i}}
{{indent}}{{end}}@{{.Name}}{{if len .Parameters}}({{range $j, $y := .Parameters}}{{if $j}}, {{end}}{{if .Name}}{{.Name}} = {{end}}{{formatValue .}}{{end}}){{end}}{{end}}
{{end}}{{end}}