// AUTO-GENERATED FILE - DO NOT MODIFY
// Generated from {{idl.Filename}}
{{setindent ""}}{{template "SIMPLECOMMENTS" idl.Comments }}
package {{package}};

import java.io.Serializable;

{{template "COMMENTS" .Comments }}
public enum {{.Name}} implements com.concur.babel.model.BabelEnum, Serializable {
{{setindent "\t"}}{{ $e := .}}{{range $i, $v := .Values}}
{{indent}}{{.Name}}({{.Value}}){{if last $i $e.Values | not}},{{else}};{{end}}{{end}}

{{indent}}private final int value;

{{indent}}public int getValue() { return this.value; }

{{indent}}private {{.Name}}(int value) {
{{indent}}{{indent}}this.value = value;
{{indent}}}

{{indent}}public static {{.Name}} findByValue(int value) {
{{indent}}{{indent}}switch(value) { {{indent}}{{indent}}{{indent}}{{range .Values}}
{{indent}}{{indent}}{{indent}}case {{.Value}}:
{{indent}}{{indent}}{{indent}}{{indent}}return {{.Name}};{{end}}
{{indent}}{{indent}}{{indent}}default:
{{indent}}{{indent}}{{indent}}{{indent}}return null;
{{indent}}{{indent}}}
{{indent}}}
}