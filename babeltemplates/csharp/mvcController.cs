// <auto-generated />
// AUTO-GENERATED FILE - DO NOT MODIFY
// Generated from {{.Filename}}

using System;
using System.Collections.Generic;
using Concur.Babel;
{{range usings}}using {{.}};
{{end}}
namespace {{index .Namespaces "csharp"}}
{ {{range $k, $s := .Services}}{{if $k}}
{{end}}
{{setindent "\t"}}{{template "COMMENTS" .Comments }}	[System.CodeDom.Compiler.GeneratedCode("Babel", "")]
	public partial class {{.Name}}Controller : {{baseController}}<I{{.Name}}>
	{ {{range .Methods}}
{{if .Parameters}}{{setindent "\t\t"}}		class {{.Name}}Request : Concur.Babel.Mvc.IBabelRequest
		{ {{range $i, $x := .Parameters}}{{setindent "\t\t\t"}}
{{template "COMMENTS" .Comments }}			{{template "ATTRS" .Attributes}}public {{formatType .Type}} {{toPascalCase .Name}};
{{end}}
			#region IBabelRequest
			public void RunOnChildren<T>(BabelModelAction<T> method, T auxData, bool runOnAll = true)
			{
{{range .Parameters}}				{{if isTrivialProperty .Type}}if(runOnAll) {{end}}{{toPascalCase .Name}} = ({{formatType .Type}}) method("{{.Name}}", typeof({{formatType .Type}}), {{toPascalCase .Name}}, auxData);
{{end}}			}

			public bool RunOnChild<T>(string name, BabelModelAction<T> method, T auxData)
			{
				switch(name)
				{
{{range .Parameters}}					case "{{.Name}}": {{toPascalCase .Name}} = ({{formatType .Type}}) method("{{.Name}}", typeof({{formatType .Type}}), {{toPascalCase .Name}}, auxData); return true;
{{end}}					default: return false;
				}
			}

			public void SetDefaults()
			{
{{range .Parameters}}{{if .Initializer}}				if ({{toPascalCase .Name}} == null) {{toPascalCase .Name}} = {{cast .Type}}{{formatValue .Initializer}};
{{end}}{{end}}			}
			#endregion
		}
{{end}}
{{setindent "\t\t"}}{{template "METHODCOMMENTS" .}}{{template "ATTRS" .Attributes}}		public {{formatType .Returns}} {{.Name}}()
		{
			{{if .Parameters}}var requestData = DeserializeRequest<{{.Name}}Request>();
			{{end}}{{if isVoid .Returns | not}}return {{end}}m_businessLogic.{{toPascalCase .Name}}({{range $i, $x := .Parameters}}{{if $i}}, {{end}}requestData.{{toPascalCase .Name}}{{end}});
		}
{{end}}	}{{end}}
}
