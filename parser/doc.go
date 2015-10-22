/*
	The IDL parser reads IDL files for the Babel tools. IDL files describe models
	and web services. The parser package uses Go's built-in lexical analyzer and
	Go's YACC tool to parse IDL files. The result is placed into the structures
	defined in the idl package.

	For more information, see the README.md file.
*/
package parser

//go:generate -command yacc go tool yacc
//go:generate yacc -o parseidl.go parseidl.y
