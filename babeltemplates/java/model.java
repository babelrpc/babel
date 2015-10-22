// AUTO-GENERATED FILE - DO NOT MODIFY
// Generated from {{idl.Filename}}
{{setindent ""}}{{template "SIMPLECOMMENTS" idl.Comments }}
package {{package}};

import com.google.gson.annotations.SerializedName;
import java.io.Serializable;
{{range imports}}import {{.}}.*;
{{end}}
{{$md := .}}
{{template "COMMENTS" .Comments }}{{template "ATTRS" .Attributes }}
public{{if .Abstract}} abstract{{end}} class {{.Name}}{{if .Extends}} extends {{.Extends}}{{end}} implements Serializable {	
{{range .Fields}}{{setindent "\t"}}
{{template "COMMENTS" .Comments }}
{{indent}}@SerializedName("{{.Name}}")
{{indent}}private {{formatType .Type}} {{toCamelCase .Name}}{{if .Initializer}} = {{formatInitializerLiteral .}}{{end}}{{if .Type.IsList}} = new {{formatListInit .Type}}(){{end}}{{if .Type.IsMap}} = new {{formatMapInit .Type}}(){{end}};{{end}}

{{setindent "\t"}}{{indent}}public {{.Name}}() {}

{{if len .Fields}}{{setindent "\t"}}{{indent}}public {{.Name}}({{range $i, $v := .Fields}}
{{indent}}{{indent}}{{formatType .Type}} {{toCamelCase .Name}}{{if last $i $.Fields | not}},{{end}}{{end}})
{{indent}}{
{{range .Fields}}
{{indent}}{{indent}}this.{{toCamelCase .Name}} = {{toCamelCase .Name}};{{end}}
{{indent}}}{{end}}
{{range .Fields}}
{{setindent "\t"}}{{indent}}public {{formatType .Type}} {{getterName .}}() { return this.{{toCamelCase .Name}}; };
{{setindent "\t"}}{{indent}}public void {{setterName .}}({{formatType .Type}} {{toCamelCase .Name}}) {
{{indent}}{{indent}}this.{{toCamelCase .Name}} = {{toCamelCase .Name}};
{{indent}}}
{{end}}	
{{setindent "\t"}}{{indent}}public String toString() {
{{if len .Fields}}{{indent}}{{indent}}StringBuilder sb = new StringBuilder("{{.Name}}(");{{ $s := .}}
{{range $i, $v := .Fields}}{{indent}}{{indent}}sb.append("{{toCamelCase .Name}}:");
{{indent}}{{indent}}sb.append(this.{{toCamelCase .Name}} + "{{if last $i $s.Fields | not}}, {{else}}){{end}}");
{{end}}{{indent}}{{indent}}return sb.toString();{{else}}
{{indent}}{{indent}}return new StringBuilder("{{.Name}}()").toString();
{{end}}
{{indent}}}
}