// AUTO-GENERATED FILE - DO NOT MODIFY
// Generated from {{idl.Filename}}
{{setindent ""}}{{template "SIMPLECOMMENTS" idl.Comments }}
package {{package}};

{{template "COMMENTS" .Comments }}
public class {{.Name}} { {{setindent "\t"}}
	
{{range .Values}}{{indent}}public static final {{constType .DataType}} {{.Name}} = {{formatValue .}};
{{end}}
	
}