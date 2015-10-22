{{define "COMMENTS"}}{{range .}}{{indent}}//{{.}}
{{end}}{{end}}
{{define "SIMPLECOMMENTS"}}{{range .}}{{indent}}//{{.}}
{{end}}{{end}}
{{define "METHODCOMMENTS"}}{{range .Comments}}{{indent}}//{{.}}
{{end}}{{range .Parameters}}{{indent}}// {{.Name}}: {{range $i,$m := .Comments}}{{if $i}}
{{indent}}// {{end}}{{.}}{{end}}
{{end}}{{end}}
