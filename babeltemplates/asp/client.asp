{{if isAsp}}<%{{end}}{{$idl := .}}
' AUTO-GENERATED FILE - DO NOT MODIFY
' Generated from {{.Filename}}
{{template "COMMENTS" $idl.Comments }}

' Requires inc_babel.asp

{{range .Services}}{{$fn := fullNameOf .Name}}{{$sn := .Name}}{{setindent "\t"}}{{template "COMMENTS" .Comments }}class {{$fn}}Client

	Private URL
	Private TimeoutSecs
	Private Headers

	sub Class_Initialize
		URL = ""
		TimeoutSecs = 30
		set Headers = Nothing
	end sub

	sub Class_Terminate
		set Headers = Nothing
	end sub

	sub SetHeader(key, value)
		if Headers is Nothing then set Headers = CreateObject("Scripting.Dictionary")
		Headers(key) = value
	end sub

	' Init
	' baseUrl - Base service URL
	' timeoutSeconds - Timeout in seconds
	sub InitHttp(baseUrl, timeoutSeconds)
		URL = baseUrl
		TimeoutSecs = timeoutSeconds
	end sub
{{range $i, $m := .Methods}}
{{setindent "\t"}}{{template "METHODCOMMENTS" .}}	{{if isVoid .Returns}}sub{{else}}function{{end}} {{.Name}}({{range $i, $v := .Parameters}}{{.Name}}{{if last $i $m.Parameters | not}}, {{else}}{{end}}{{end}})
		dim parms : set parms = CreateObject("Scripting.Dictionary")
{{range $i, $v := .Parameters}}{{indent}}{{indent}}{{if .Type.IsStruct $idl}}set {{end}}{{if .Type.IsMap}}set {{end}}parms("{{.Name}}") = {{.Name}} ' {{.Type}}
{{end}}		{{if isVoid .Returns}}Call {{else}}{{if or (.Returns.IsStruct $idl) .Returns.IsMap}}set {{end}}{{.Name}} = {{end}}BabelCall(URL, "{{$sn}}/{{.Name}}", Headers, TimeoutSecs, parms, "{{internalType .Returns}}")
	end {{if isVoid .Returns}}sub{{else}}function{{end}}
{{end}}
end class

{{end}}{{if isAsp}}%>{{end}}