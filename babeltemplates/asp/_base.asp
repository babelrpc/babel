{{define "COMMENTS"}}{{$cmts := expandComments .}}{{range $cmts}}{{indent}}' {{.}}
{{end}}{{end}}{{define "ATTRS"}}{{$attrs := filterAttrs .}}{{if len $attrs}}{{indent}}'[{{range $i, $x := $attrs}}{{if $i}}, {{end}}{{.Name}}{{if len .Parameters}}({{range $j, $y := .Parameters}}{{if $j}}, {{end}}{{if .Name}}{{.Name}} = {{end}}{{formatValue .}}{{end}}){{end}}{{end}}]
{{end}}{{end}}{{define "METHODCOMMENTS"}}{{$cmts := expandComments .Comments}}{{range $cmts}}{{indent}}' {{.}}
{{end}}{{range .Parameters}}{{indent}}' {{.Name}} - {{range $i,$m := .Comments}}{{.}}{{end}}
{{end}}{{end}}
