{{define "COMMENTS"}}{{$cmts := expandComments .}}{{if len $cmts}}{{indent}}/// <summary>
{{range .}}{{indent}}/// {{.}}
{{end}}{{indent}}/// </summary>
{{end}}{{end}}
{{define "SIMPLECOMMENTS"}}{{range .}}{{indent}}//{{.}}
{{end}}{{end}}
{{define "METHODCOMMENTS"}}{{$cmts := expandComments .Comments}}{{indent}}/// <summary>
{{range $cmts}}{{indent}}/// {{.}}
{{end}}{{indent}}/// </summary>
{{range .Parameters}}{{indent}}/// <param name="{{.Name}}">{{range $i,$m := .Comments}}{{if $i}}
{{indent}}/// {{end}}{{.}}{{end}}</param>
{{end}}{{end}}
{{define "ATTRS"}}{{$attrs := filterAttrs .}}{{if len $attrs}}{{indent}}[{{range $i, $x := $attrs}}{{if $i}}, {{end}}{{.Name}}{{if len .Parameters}}({{range $j, $y := .Parameters}}{{if $j}}, {{end}}{{if .Name}}{{.Name}} = {{end}}{{formatValue .}}{{end}}){{end}}{{end}}]
{{end}}{{end}}