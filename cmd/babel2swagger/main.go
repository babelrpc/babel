/*
	babel2swagger converts Babel IDL to Swagger 2.0 format.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/babelrpc/babel/idl"
	"github.com/babelrpc/babel/parser"
	"github.com/babelrpc/swagger2"
	"html"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// This attached swagger for the error file onto this binary, to use later
//go:generate babel2swagger -out error.json ../../babeltemplates/error.babel
//go:generate binder -o errormodel.go error.json

// Flags that are global
var (
	flatten bool
	restful bool
	swagInt bool
	genErr  bool
)

// main entry point
func main() {
	/*
		// recover panicking parser
		defer func() {
			if r := recover(); r != nil {
				e, y := r.(*idl.Error)
				if y {
					fmt.Fprintln(os.Stderr, e.Error())
					os.Exit(e.Code)
				} else {
					fmt.Fprintf(os.Stderr, "Exiting due to fatal error:\n%s\n", r)
					os.Exit(1)
				}
			}
		}()
	*/

	var format string
	flag.StringVar(&format, "format", "json", "Specifies output format - can be json or yaml")

	var output string
	flag.StringVar(&output, "out", "", "Specifies the file to write to")

	var host string
	flag.StringVar(&host, "host", "localhost", "Specifies the host to include in the file, for example localhost:8080")

	var basePath string
	flag.StringVar(&basePath, "basepath", "/", "Specifies the base path to include in the file, for example /foo/bar")

	var title string
	flag.StringVar(&title, "title", "My Application", "Sets the application title")

	flag.BoolVar(&flatten, "flat", false, "Flatten composed objects into a single object definition")

	var version string
	flag.StringVar(&version, "version", "1.0", "Sets the API version")

	flag.BoolVar(&restful, "rest", false, "Process @rest annotations (resulting Swagger won't be able to invoke Babel services)")

	flag.BoolVar(&swagInt, "int64", false, "When -rest is enabled, format int64 Swagger-style instead of Babel-style")

	flag.BoolVar(&genErr, "error", false, "When -rest is enabled, still include the Babel error definition")

	flag.Parse()

	if format != "json" && format != "yaml" {
		fmt.Printf("-format must be json or yaml.\n")
		os.Exit(0)
	}

	if swagInt && !restful {
		fmt.Fprintf(os.Stderr, "Warning: Ignoring -int64 because -rest is not set.\n")
		swagInt = false
	}

	if len(flag.Args()) == 0 {
		fmt.Printf("Please specify files to process or -help for options.\n")
		os.Exit(0)
	}

	// initialize map to track processed files
	processedFiles := make(map[string]bool)

	// create base IDL to aggregate into
	var midl idl.Idl
	midl.Init()
	midl.AddDefaultNamespace("babelrpc.io", "Foo/Bar")

	// load specified IDL files into it
	for _, infilePat := range flag.Args() {
		infiles, err := filepath.Glob(infilePat)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot glob files: %s\n", err)
			os.Exit(5)
		}
		if len(infiles) == 0 {
			fmt.Fprintf(os.Stderr, "Warning: No files match \"%s\"\n", infilePat)
		}
		for _, infile := range infiles {
			_, ok := processedFiles[infile]
			if ok {
				fmt.Fprintf(os.Stderr, "Already processed %s\n", infile)
			} else {
				processedFiles[infile] = true
				// fmt.Printf("%s:\n", infile)
				bidl, err := parser.ParseIdl(infile, "test")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Parsing error in %s: %s\n", infile, err)
					os.Exit(6)
				}

				midl.Imports = append(midl.Imports, bidl)
			}
		}
	}

	// validate combined babel
	err := midl.Validate("test")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Combined IDL does not validate: %s\n", err)
	}

	// convert to swagger

	var swag swagger2.Swagger
	swag.Swagger = "2.0"
	swag.Host = host
	swag.BasePath = basePath
	swag.Consumes = []string{"application/json"}
	swag.Produces = []string{"application/json"}
	swag.Tags = []swagger2.Tag{swagger2.Tag{Name: "babel"}}
	swag.Info.Title = title
	swag.Info.Version = version
	swag.Definitions = make(swagger2.Definitions, 0)

	// parse and attach error model (not included with REST unless -genErr is used)
	if !restful || genErr {
		errorB := Lookup("/error.json")
		errorS, err := swagger2.LoadJson(errorB)
		if err != nil {
			panic("Cannot parse embedded error file")
		}
		// attach error definitions
		for sx, sy := range errorS.Definitions {
			swag.Definitions[sx] = sy
		}
	}

	// Glob together file-level commebts
	sarr := make([]string, 0)
	for _, idl := range midl.Imports {
		if len(idl.Comments) > 0 {
			sarr = append(sarr, idl.Filename+":\n\n")
			sarr = append(sarr, idl.Comments...)
		}
	}
	swag.Info.Description = strings.Join(sarr, "\n\n")
	if swag.Info.Description != "" {
		swag.Info.Description += "\n\n"
	}
	gentxt := "Generated on " + time.Now().Format(time.RFC1123) + " via:\n\n\tbabel2swagger"
	if format != "json" {
		gentxt += " -format " + format
	}
	if flatten {
		gentxt += " -flat"
	}
	if restful {
		gentxt += " -rest"
		if swagInt {
			gentxt += " -int64"
		}
		if genErr {
			gentxt += " -error"
		}
	}
	if host != "localhost" {
		gentxt += " -host " + host
	}
	if basePath != "/" {
		gentxt += " -basepath " + basePath
	}
	if version != "1.0" {
		gentxt += " -version " + version
	}
	if title != "My Application" {
		gentxt += " -title \"" + title + "\""
	}
	for _, g := range flag.Args() {
		gentxt += " \"" + g + "\""
	}
	swag.Info.Description += gentxt

	// add structs to swagger definitions
	for _, st := range allStructs(&midl) {
		swag.Definitions[st.Name] = *structToSchema(&midl, st)
	}

	// add enums to swagger definitions
	// SWAGGER-BUG: swagger-ui doesn't like it this way - need expand them out individually as strings
	/*
		for _, en := range allEnums(&midl) {
			swag.Definitions[en.Name] = *enumToSchema(&midl, en)
		}
	*/

	// track whether we find any empty parameters
	hasEmptyParms := false

	// add service/methods to paths
	swag.Paths = make(swagger2.Paths, 0)
	for _, svc := range allServices(&midl) {
		if restful {
			err := addRestService(&swag, &midl, svc)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Cannot generate REST service: %s\n", err)
				os.Exit(11)
			}
		} else {
			svcComments := strings.Join(svc.Comments, "\n")
			/*
				Seems like these belong somewhere else
				if svcComments != "" {
					swag.Info.Description += "\n" + svcComments
				}
			*/
			for _, mth := range svc.Methods {
				var p swagger2.PathItem
				p.Post = new(swagger2.Operation)
				p.Post.Tags = []string{svc.Name}
				mthComments := strings.Join(mth.Comments, "\n")
				theseArgs := make([]string, 0)
				for _, s := range mth.Parameters {
					// SWAGGER-BUG: Swagger should be HTML escaping this
					theseArgs = append(theseArgs, html.EscapeString(s.Type.String())+" "+s.Name)
				}
				// SWAGGER-BUG: Swagger should be HTML escaping this
				p.Post.Description = html.EscapeString(mth.Returns.String()) + " " + mth.Name + "(" + strings.Join(theseArgs, ", ") + ")"
				if mthComments != "" {
					p.Post.Description += "\n\n" + mthComments
				}
				if svcComments != "" {
					p.Post.Description += "\n\n" + svc.Name + ": " + svcComments
				}
				p.Post.OperationId = svc.Name + "_" + mth.Name
				p.Post.Summary = svc.Name + "." + mth.Name
				p.Post.Parameters = make([]swagger2.Parameter, 0)
				var parm swagger2.Parameter
				parm.Name = "request"
				parm.In = "body"
				parm.Required = new(bool)
				if mth.HasParameters() {
					oname := strings.Title(svc.Name) + strings.Title(mth.Name) + "RequestArgs"
					swag.Definitions[oname] = *parmsToSchema(&midl, mth)
					parm.Schema = &swagger2.Schema{ItemsDef: swagger2.ItemsDef{Ref: "#/definitions/" + oname}}
					parm.Description = "Request definition"
					*parm.Required = true
				} else {
					parm.Schema = &swagger2.Schema{ItemsDef: swagger2.ItemsDef{Ref: "#/definitions/EmptyRequestArgs"}}
					hasEmptyParms = true
					*parm.Required = false
				}
				p.Post.Parameters = append(p.Post.Parameters, parm)
				p.Post.Responses = make(swagger2.Responses, 0)
				// SWAGGER-BUG: note that swagger-ui does not show primitive types for responses, even though they are allowed.
				p.Post.Responses["200"] = swagger2.Response{
					Description: "Response of type " + html.EscapeString(mth.Returns.String()),
					Schema:      returnsToSchema(&midl, mth.Returns),
				}
				p.Post.Responses["default"] = swagger2.Response{
					Description: "error",
					Schema:      &swagger2.Schema{ItemsDef: swagger2.ItemsDef{Ref: "#/definitions/ServiceError"}},
				}
				swag.Paths["/"+svc.Name+"/"+mth.Name] = p
			}
		}
	}

	// Add empty parameter definition if needed
	if hasEmptyParms {
		swag.Definitions["EmptyRequestArgs"] = swagger2.Schema{ItemsDef: swagger2.ItemsDef{Type: "object"}, Description: "Empty object"}
	}

	// Validate swagger
	errs := swag.Validate()
	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "Warning: Swagger does not validate\n%s\n", swagger2.ErrorList(errs).Indent("\t"))
	}

	// open output
	var outfile io.WriteCloser
	outfile = os.Stdout
	if output != "" {
		outfile, err = os.Create(output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Cannot open output file %s\n", output)
			os.Exit(9)
		}
		defer outfile.Close()
	}

	// Write output
	if format == "json" {
		s, err := swag.Json()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot convert to JSON: %s\n", err)
			os.Exit(7)
		}
		fmt.Fprintln(outfile, string(s))
	} else {
		s, err := swag.Yaml()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot convert to YAML: %s\n", err)
			os.Exit(8)
		}
		fmt.Fprintln(outfile, string(s))
	}
}
