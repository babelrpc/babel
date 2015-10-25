![Babel](http://babelrpc.io/media/logo.png)

Babel is an IDL parser and RPC framework using JSON over HTTP. IDL files describe models and web services. The `babel` tool allows you to generate client and server code in multiple languages from the IDL file.

Visit the [Babel RPC Home Page](http://babelrpc.io) for more information.

## Babel Tools

The babel tools are:

* [allbabeltypes](cmd/allbabeltypes) - A test tool that generates a babel file containing most possible combinations of types, for testing.
* [babel](cmd/babel) - The [IDL](idl) compiler.
* [babel2swagger](cmd/babel2swagger) - A tool to convert Babel to Swagger 2.
* [babelproxy](cmd/babelproxy) - A tool to use [rest annotations](rest) to proxy RESTful APIs for a babel service.

The main Babel libraries are:

* [babeltemplates](babeltemplates) - Language templates for Babel.
* [generator](generator) - Code for language-specific code generators.
* [idl](idl) - Code for Babel's Interface Definition Language.
* [parser](parser) - A [YACC](http://golang.org/cmd/yacc/)-based parser for Babel files.
* [rest](rest) - Process RESTful annotations (attributes) in Babel files.
* [swagger2](https://github.com/babelrpc/swagger2) - Serialize Swagger 2 structures to JSON and YAML.
