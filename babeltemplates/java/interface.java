// AUTO-GENERATED FILE - DO NOT MODIFY
// Generated from {{idl.Filename}}
{{setindent ""}}{{template "SIMPLECOMMENTS" idl.Comments }}
package {{package}};

import java.util.HashMap;
import java.util.Map;

import com.concur.babel.ServiceMethod;
import com.concur.babel.BabelService;
import com.concur.babel.ResponseServiceMethod;
import com.concur.babel.VoidServiceMethod;
import com.concur.babel.processor.BaseInvoker;
import com.concur.babel.transport.BaseClient;
import com.concur.babel.transport.Transport;
import com.concur.babel.service.BabelServiceDefinition;
{{$srv := .}}
{{range imports}}import {{.}}.*;
{{end}}
{{template "COMMENTS" .Comments }}
public class {{.Name}} implements BabelServiceDefinition { {{setindent "\t"}}

{{indent}}/**
{{indent}} * Gets the interface class for this service.
{{indent}} */
{{indent}}public Class<Iface> getIfaceClass() { return Iface.class; }

{{indent}}/**
{{indent}} * Convenience method to create a new Invoker instance containing the service iface implementation.
{{indent}} */
{{indent}}public Invoker createInvoker(BabelService iFaceImpl) { return new Invoker((Iface)iFaceImpl); }

{{indent}}/**
{{indent}} * The interface defining the methods for this service. You should provide an implmentation of this interface. ServiceName.Iface.
{{indent}} */
{{indent}}public interface Iface extends BabelService {
{{indent}}{{indent}}{{range $i, $m := .Methods}}{{template "COMMENTS" .Comments }}
{{indent}}{{indent}}{{formatType .Returns}} {{toCamelCase .Name}}({{range $i, $v := .Parameters}}{{formatType .Type}} {{toCamelCase .Name}}{{if last $i $m.Parameters | not}}, {{else}}{{end}}{{end}});
{{end}}
{{indent}}}
{{indent}}/**
{{indent}} * Client can be created to make a service call for any of the methods defined in this service.
{{indent}} */
{{indent}}public static class Client extends BaseClient implements Iface {
	
{{indent}}{{indent}}public Client(String url) { super(url); }
{{indent}}{{indent}}public Client(String url, int timeoutInMillis) { super(url, timeoutInMillis); }
{{indent}}{{indent}}public Client(Transport transport) { super(transport); }
{{indent}}{{indent}}{{range $i, $m := .Methods}}
{{indent}}{{indent}}public {{formatType .Returns}} {{toCamelCase .Name}}({{range $i, $v := .Parameters}}{{formatType .Type}} {{toCamelCase .Name}}{{if last $i $m.Parameters | not}}, {{else}}{{end}}{{end}}) {
	
{{indent}}{{indent}}{{indent}}{{toPascalCase .Name}} serviceMethod = new {{toPascalCase .Name}}({{range $i, $v := .Parameters}}{{toCamelCase .Name}}{{if last $i $m.Parameters | not}}, {{else}}{{end}}{{end}});
{{indent}}{{indent}}{{indent}}{{if isVoid .Returns}}{{else}}return {{end}}this.transport.invoke(serviceMethod);

{{indent}}{{indent}}}
{{end}}
{{indent}}}

{{indent}}/**
{{indent}} * Invoker is a component in the babel framework used to call service methods.  This component contains the service iface implementation and needs to be registered with the ServiceRequestDispatcher.
{{indent}} */
{{indent}}public static class Invoker extends BaseInvoker<Iface> {
	
{{indent}}{{indent}}public Invoker(Iface serviceImpl) { super(serviceImpl); }

{{indent}}{{indent}}public Map<String, Class<? extends ServiceMethod>> initServiceMethods() {
{{indent}}{{indent}}{{indent}}Map<String, Class<? extends ServiceMethod>> map = new HashMap<String, Class<? extends ServiceMethod>>();
{{indent}}{{indent}}{{indent}}{{range $i, $m := .Methods}}map.put("{{toCamelCase .Name}}", {{toPascalCase .Name}}.class);{{end}}
{{indent}}{{indent}}{{indent}}return map;
{{indent}}{{indent}}}

{{indent}}{{indent}}public String getServiceName() { return "{{$srv.Name}}"; }
{{indent}}{{indent}}public Class<Iface> getInterface() { return Iface.class; }

{{indent}}}

{{range $i, $m := .Methods}}
{{indent}}private static class {{toPascalCase .Name}} extends {{if isVoid .Returns}}VoidServiceMethod{{else}}ResponseServiceMethod<{{parseType .Returns}}>{{end}} {	
{{range $i, $v := .Parameters}}
{{indent}}{{indent}}private {{formatType .Type}} {{toCamelCase .Name}}{{if .Initializer}} = {{formatInitializerLiteral .}}{{end}};{{end}}

{{if len .Parameters}}{{indent}}{{indent}}public {{toPascalCase .Name}}() {}{{end}}
{{indent}}{{indent}}public {{toPascalCase .Name}}({{range $i, $v := .Parameters}}{{formatType .Type}} {{toCamelCase .Name}}{{if last $i $m.Parameters | not}}, {{else}}{{end}}{{end}}) {
	
{{range $i, $v := .Parameters}}{{indent}}{{indent}}{{indent}}this.{{toCamelCase .Name}} = {{toCamelCase .Name}};
{{end}}{{indent}}{{indent}}}

{{indent}}{{indent}}public String getServiceName() { return "{{$srv.Name}}"; }
{{indent}}{{indent}}public String getMethodName() { return "{{toCamelCase .Name}}"; }

{{indent}}{{indent}}public Object[] getMethodParameters() {
{{indent}}{{indent}}{{indent}}return new Object[] { {{range $i, $v := .Parameters}}this.{{toCamelCase .Name}}{{if last $i $m.Parameters | not}}, {{else}}{{end}}{{end}} };
{{indent}}{{indent}}}
{{indent}}}
{{end}}
}