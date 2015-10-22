{{define "INITFIELDS"}}{{if len .Fields}}
		' Fields from {{.Name}}{{end}}{{range .Fields}}{{if .Type.IsStruct idl}}
		set {{.Name}} = Nothing{{end}}{{if .Type.IsMap}}
		set {{.Name}} = CreateObject("Scripting.Dictionary"){{end}}{{if .Type.IsList}}
		{{.Name}} = Array(){{end}}{{if .Initializer}}
		{{.Name}} = {{formatValue .Initializer}}{{end}}{{end}}{{end}}{{define "CLOSEFIELDS"}}{{if len .Fields}}
		' Fields from {{.Name}}{{end}}{{range .Fields}}{{if .Type.IsStruct idl}}
		set {{.Name}} = Nothing{{end}}{{if .Type.IsMap}}
		set {{.Name}} = Nothing{{end}}{{if .Type.IsList}}
		{{.Name}} = Empty{{end}}{{end}}{{end}}{{define "FIELDS"}}{{range .Fields}}
{{setindent "\t"}}{{template "COMMENTS" .Comments }}{{template "ATTRS" .Attributes}}	public {{.Name}} ' {{.Type}}{{if (.Type.IsEnum idl)}} - see "{{fullNameOf .Type.Name}}" for values{{end}}
{{end}}{{end}}{{define "TOJSON"}}{{setindent "\t"}}{{range $i, $v := .Fields}}
{{indent}}{{indent}}call s_.Write(j_, "{{internalType $v.Type}}", "{{$v.Name}}", {{$v.Name}}, "{{renames $v.Type}}", i_){{end}}{{end}}{{define "FROMJSON"}}{{setindent "\t"}}{{range $i, $v := .Fields}}
{{indent}}{{indent}}{{if or ($v.Type.IsStruct idl) $v.Type.IsMap}}set {{end}}{{$v.Name}} = s_.Read(j_, "{{internalType $v.Type}}", "{{$v.Name}}", {{$v.Name}}, "{{renames $v.Type}}"){{end}}{{end}}{{if isAsp}}<%{{end}}{{$idl := .}}
' AUTO-GENERATED FILE - DO NOT MODIFY
' Generated from {{.Filename}}
{{template "COMMENTS" $idl.Comments }}

{{if len .Enums}}' *** Enumerations ***

{{end}}{{range .Enums}}{{$fn := fullNameOf .Name}}{{setindent ""}}{{template "COMMENTS" .Comments }}{{range .Values}}const {{$fn}}{{.Name}} = "{{.Name}}" ' Value of {{formatValue .}}
{{end}}
function {{$fn}}_ToId(strVal)
	dim r
	select case strVal
{{range .Values}}	case "{{.Name}}"
		r = {{formatValue .}}
{{end}}	case else
		r = empty
	end select
	{{$fn}}_ToId = r
end function

function {{$fn}}_ToString(idVal)
	dim r
	select case idVal
{{range .Values}}	case "{{formatValue .}}"
		r = "{{.Name}}"
{{end}}	case else
		r = empty
	end select
	{{$fn}}_ToString = r
end function
{{end}}
{{if len .Consts}}' *** Constants ***

{{end}}{{range .Consts}}{{$fn := fullNameOf .Name}}{{setindent ""}}{{template "COMMENTS" .Comments }}{{range .Values}}const {{$fn}}{{.Name}} = {{formatValue .}}
{{end}}
{{end}}
{{if len .Structs}}' *** Structures ***

{{end}}{{range .Structs}}{{$fn := fullNameOf .Name}}{{$bases := .BaseClasses $idl}}{{$subs := .SubClasses $idl}}{{setindent ""}}{{template "COMMENTS" .Comments }}{{template "ATTRS" .Attributes}}class {{$fn}}{{if len $bases}}
	' inheritance: {{range $i, $x := $bases}}{{if $i}} -> {{end}}{{if .Abstract}}abstract {{end}}{{fullNameOf .Name}}{{end}} -> {{fullNameOf .Name}}{{end}}{{if len $subs}}
	' Subclasses:{{range $subs}} {{.Name}}{{end}}{{end}}

	sub Class_Initialize{{range $i, $x := $bases}}{{template "INITFIELDS" .}}{{end}}{{template "INITFIELDS" .}}
	end sub

	sub Class_Terminate{{range $i, $x := $bases}}{{template "CLOSEFIELDS" .}}{{end}}{{template "CLOSEFIELDS" .}}
	end sub

{{range $i, $x := $bases}}	' --- Members from {{fullNameOf .Name}} ---
{{template "FIELDS" .}}
{{end}}	' --- Members ---
{{template "FIELDS" .}}
	' Called by Babel protocol to write this object
	public sub Write(s_, j_)
		dim i_ : i_ = 0{{range $i, $x := $bases}}{{template "TOJSON" .}}{{end}}{{template "TOJSON" .}}
	end sub

	' Called by Babel protocol to read this object
	public sub Read(s_, j_){{range $i, $x := $bases}}{{template "FROMJSON" .}}{{end}}{{template "FROMJSON" .}}
	end sub

	' Convert this object to a JSON string
	public function ToJSON()
		ToJSON = BabelToJson(empty, empty, me)
	end function

	' Convert this object to an XML document
	public function ToXML()
		set ToXML = BabelToXml(empty, "{{.Name}}", me)
	end function
end class

{{end}}{{if isAsp}}%>{{end}}