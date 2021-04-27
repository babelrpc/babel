/*
	Babel is an IDL parser and code generator for web services. IDL files describe models
	and web services. The babel tool allows you to generate client and server code
	in multiple languages from the IDL file.

	For more information, see the README.md file.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/babelrpc/babel/generator"
	"github.com/babelrpc/babel/idl"
	"github.com/babelrpc/babel/parser"
)

// main entry point
func main() {
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

	templatesDir := flag.String("templates", generator.LocateTemplateDir(), "Overrides the location of the templates folder")
	outputJson := flag.Bool("json", false, "Output parse tree in JSON format")
	lang := flag.String("lang", "", "Generate code with given language csharp|java")
	outputDir := flag.String("output", "", "Output folder for generates files, defaults to gen-<<lang>>")
	inc := flag.Bool("inc", false, "Generate included files too")
	getHelp := flag.Bool("help", false, "Get command-line help")
	scopes := flag.String("scopes", "", "Comma-separated list of scopes to enable (no spaces)")
	options := flag.String("options", "", "Comma-separated list of key=value pairs for generator")
	genClient := flag.Bool("client", false, "Generate client files")
	genServer := flag.Bool("server", false, "Generate server files")
	genModel := flag.Bool("model", false, "Generate model files")
	serverType := flag.String("servertype", "", "Optional language-specific server type")
	ver := flag.Bool("version", false, "Display Babel version number")
	nsMatch := flag.String("ns", "", "Optionally matches files only if namespace starts with this")
	flag.Parse()

	if *getHelp {
		fmt.Printf("The babel command generates source files from Babel IDL files.\n\n")
		fmt.Printf("babel -lang <language> [optional flags] <filePattern> [filePattern...]\n\n")

		flag.PrintDefaults()

		fmt.Printf(`
If none of -model, -client, or -server are specified, then all three are generated.

Use -scopes to enable attributes that are qualified with a scope.

-options are values that are specific to each language. See the documenation for more information.
	ASP supports "ext", which can be "vbs" or "asp".
	C# supports "controller", which can be used to override the controller base class.
	C# supports "output", which can be used to define output directory options.  Supported options are:
		ns-flat - code is generated in a single (flat) output directory corresponding to the namespace.
		ns-nested - code is generated in a multiple (nested) output directories corresponding to the namespace.

Babel accepts quite flexible file patterns - and accepts more than one. Here are some examples:

	*.babel         All babel files in the current folder
	*/*.babel       All babel files in every immediate folder
	foo/*/*.babel   All babel files in each folder of foo

`)

		os.Exit(0)
	}

	if *ver {
		fmt.Printf("Babel Version %s\n", BABEL_VERSION)
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		fmt.Printf("Please specify files to process or -help for options.\n")
		os.Exit(0)
	}

	if strings.TrimSpace(*lang) == "" {
		fmt.Fprintf(os.Stderr, "-lang must be specified\n")
		os.Exit(1)
	}

	if !*genClient && !*genServer && !*genModel {
		*genClient = true
		*genServer = true
		*genModel = true
	}

	if *outputDir == "" {
		*outputDir = "gen-" + *lang
	}

	if *outputDir != "" {
		err := os.MkdirAll(*outputDir, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output directory: %s\n", err)
			os.Exit(2)
		}
	}

	theScopes := make([]string, 0)
	for _, s := range strings.Split(*scopes, ",") {
		sc := strings.TrimSpace(s)
		if sc != "" {
			theScopes = append(theScopes, sc)
		}
	}

	theOptions := make(map[string]string)
	for _, s := range strings.Split(*options, ",") {
		sv := strings.TrimSpace(s)
		if sv != "" {
			kv := strings.Split(sv, "=")
			if len(kv) != 2 || strings.TrimSpace(kv[0]) == "" {
				fmt.Fprintf(os.Stderr, "Invalid options: %s\n", s)
				os.Exit(3)
			}
			theOptions[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}

	gen, err := generator.New(*lang, &generator.Arguments{
		OutputDir:   *outputDir,
		TemplateDir: *templatesDir,
		Scopes:      theScopes,
		GenServer:   *genServer,
		GenClient:   *genClient,
		GenModel:    *genModel,
		Options:     theOptions,
		ServerType:  *serverType})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating %s generator: %s\n", *lang, err)
		os.Exit(4)
	}

	processedFiles := make(map[string]bool)
	generatedFiles := make(map[string]bool)

	for _, infilePat := range flag.Args() {
		infiles, err := filepath.Glob(infilePat)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot glob files:\n%s\n", err)
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
				fmt.Printf("%s:\n", infile)
				bidl, err := parser.ParseIdl(infile, *lang)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Parsing error:\n%s\n", err)
					os.Exit(6)
				}

				if strings.HasPrefix(bidl.Namespaces["#default"], *nsMatch) {
					if *outputJson {
						b, err := json.MarshalIndent(bidl, "", "  ")
						if err != nil {
							fmt.Fprintf(os.Stderr, "json error: %s\n", err)
						}
						fmt.Println(string(b))
						continue
					}

					files, err := gen.GenerateCode(bidl)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error generating files for %s: %s\n", filepath.FromSlash(bidl.Filename), err)
						os.Exit(7)
					}
					fmt.Println("\t" + strings.Join(files, "\n\t"))
					for _, gfn := range files {
						_, ok := generatedFiles[gfn]
						if ok {
							fmt.Fprintf(os.Stderr, "Error - file %s generated in a prior step and overwritten while processing %s\n", gfn, filepath.FromSlash(bidl.Filename))
							os.Exit(8)
						} else {
							generatedFiles[gfn] = true
						}
					}
				} else {
					fmt.Println("\tSkipped.")
				}

				if *inc {
					for _, imp := range bidl.UniqueImports() {
						impFn := filepath.FromSlash(imp.Filename)
						_, ok := processedFiles[impFn]
						if ok {
							fmt.Fprintf(os.Stderr, "Already processed %s\n", impFn)
						} else {
							processedFiles[impFn] = true
							fmt.Printf("%s:\n", impFn)

							if strings.HasPrefix(imp.Namespaces["#default"], *nsMatch) {
								ifiles, err := gen.GenerateCode(imp)
								if err != nil {
									fmt.Fprintf(os.Stderr, "Error generating files for %s: %s\n", impFn, err)
									os.Exit(9)
								}
								fmt.Println("\t" + strings.Join(ifiles, "\n\t"))
								for _, gfn := range ifiles {
									_, ok := generatedFiles[gfn]
									if ok {
										fmt.Fprintf(os.Stderr, "Error - file %s generated in a prior step and overwritten while processing %s\n", gfn, impFn)
										os.Exit(10)
									} else {
										generatedFiles[gfn] = true
									}
								}
							} else {
								fmt.Println("\tSkipped.")
							}
						}
					}
				}
			}
		}
	}
}
