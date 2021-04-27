package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/babelrpc/babel/idl"
)

var (
	// Mapping of IDL types to jsscript types.
	jsTypes = map[string]string{
		"bool":     "var",
		"byte":     "var",
		"int8":     "var",
		"int16":    "var",
		"int32":    "var",
		"int64":    "var",
		"float32":  "var",
		"float64":  "var",
		"string":   "var",
		"datetime": "var",
		"decimal":  "var",
		"char":     "var",
		"binary":   "[]",
		"list":     "[]",
		"map":      "{}",
	}

	// Mapping of constants to their jsscript type
	jsConstTypes = map[string]string{
		"int":    "var",
		"float":  "var",
		"string": "var",
		"bool":   "var",
		"char":   "var",
	}
)

// jsGenerator is the code generator for jsscript.
type jsGenerator struct {
	templateManager
	args               *Arguments
	baseControllerName string
}

// formatType returns the Type as a string in format suitable for jsscript.
func (gen *jsGenerator) formatType(t *idl.Type) string {
	var s string
	ms, ok := jsTypes[t.Name]
	if !ok {
		ms = t.Name
	}
	if t.Name == "list" {
		s = fmt.Sprintf(ms, gen.formatType(t.ValueType))
	} else if t.Name == "map" {
		s = fmt.Sprintf(ms, jsTypes[t.KeyType.Name], gen.formatType(t.ValueType))
	} else if t.IsPrimitive() && t.Name != "string" {
		s = ms + "?"
	} else if t.IsEnum(gen.tplRootIdl) {
		s = ms + "?"
	} else {
		s = ms
	}
	return s
}

// fullTypeName returns the fully qualified type name using namespace syntax for jsscript.
// Note: this probably doesn't work.
func (gen *jsGenerator) fullTypeName(t *idl.Type) string {
	var s string
	ms, ok := jsTypes[t.Name]
	if !ok {
		ms = t.Name
	}
	if t.Name == "list" {
		s = fmt.Sprintf(ms, gen.fullTypeName(t.ValueType))
	} else if t.Name == "map" {
		s = fmt.Sprintf(ms, jsTypes[t.KeyType.Name], gen.fullTypeName(t.ValueType))
	} else {
		ns := gen.tplRootIdl.NamespaceOf(ms, "js")
		if ns != "" {
			s = fmt.Sprintf("%s.%s", ns, ms)
		} else {
			s = ms
		}
	}
	return s
}

// fullNameOf returns the fully qualified name of the object with its
// namespace formatted for jsscript.
func (gen *jsGenerator) fullNameOf(name string) string {
	arr := strings.Split(name, ".")
	if len(arr) > 0 {
		ns := gen.tplRootIdl.NamespaceOf(arr[0], "js")
		if ns != "" {
			return ns + "." + strings.Join(arr, ".")
		}
	}
	return name
}

// formatLiteral formats a literal value for the jsscript parser.
func (gen *jsGenerator) formatLiteral(value interface{}, typeName string) string {
	if typeName == "#ref" {
		return gen.fullNameOf(value.(string))
	} else if typeName == "char" {
		return fmt.Sprintf("%q", value)
	} else {
		return fmt.Sprintf("%#v", value)
	}
}

// filterAttributes returns the attributes for the enabled scopes.
func (gen *jsGenerator) filterAttributes(attrs []*idl.Attribute) []*idl.Attribute {
	result := make([]*idl.Attribute, 0)
	for _, s := range gen.args.Scopes {
		cs := strings.ToLower(s)
		for _, a := range attrs {
			if strings.ToLower(a.Scope) == cs {
				result = append(result, a)
			}
		}
	}
	return result
}

func (gen *jsGenerator) isVoid(t *idl.Type) bool {
	return t.IsVoid()
}

func (gen *jsGenerator) isTrivialProperty(t *idl.Type) bool {
	return t.IsPrimitive() || t.Name == "binary" || t.IsEnum(gen.tplRootIdl)
}

// init sets up the generator for use and loads the templates.
func (gen *jsGenerator) init(args *Arguments) error {
	if !args.GenClient && !args.GenModel && !args.GenServer {
		return fmt.Errorf("nothing to do")
	} else if len(args.Options) > 0 {
		for k, v := range args.Options {
			switch k {
			case "controller":
				if v == "" {
					return fmt.Errorf("controller cannot be empty")
				}
				gen.baseControllerName = v
			case "output":
				if v == "" {
					return fmt.Errorf("output option cannot be empty.  Valid options are 'ns-flat' and 'ns-nested'")
				} else if v != "ns-flat" && v != "ns-nested" {
					return fmt.Errorf("invalid output option: %s: valid options are 'ns-flat' and 'ns-nested'", v)
				}
			default:
				return fmt.Errorf("the %s option is not applicable to language js", k)
			}
		}
	}
	gen.args = args
	return gen.loadTempates(args.TemplateDir, "js", template.FuncMap{
		"formatType":        func(t *idl.Type) string { return gen.formatType(t) },
		"fullNameOf":        func(name string) string { return gen.fullNameOf(name) },
		"formatValue":       func(p *idl.Pair) string { return gen.formatLiteral(p.Value, p.DataType) },
		"filterAttrs":       func(attrs []*idl.Attribute) []*idl.Attribute { return gen.filterAttributes(attrs) },
		"isVoid":            func(t *idl.Type) bool { return gen.isVoid(t) },
		"isTrivialProperty": func(t *idl.Type) bool { return gen.isTrivialProperty(t) },
		"usings": func() []string {
			pkg := gen.tplRootIdl.Namespaces["js"]
			imports := make([]string, 0)
			for _, i := range gen.tplRootIdl.UniqueNamespaces("js") {
				if i != pkg {
					//relPath := i[len(pkg)];
					//fmt.Println(relPath);
					//fmt.Println("pkg => " + pkg);
					//fmt.Println("i => " + i);
					relPath := ""
					for x := 0; x < len(strings.Split(pkg, ".")); x++ {
						relPath += "../"
					}
					//fmt.Println("relPath => " + relPath);
					imports = append(imports, relPath+strings.Replace(i, ".", "/", -1))
				}
			}
			return imports
		},
		"isNotPascalCase": func(name string) bool {
			if len(name) > 1 {
				return strings.ToUpper(name[0:1]) != name[0:1]
			}
			return false
		},
		"baseController": func() string {
			if gen.baseControllerName != "" {
				return gen.baseControllerName
			} else {
				return "Concur.Babel.Mvc.BabelController"
			}
		},
		"cast": func(t *idl.Type) string {
			if t.Name == "float32" {
				return "(float)"
			} else {
				return ""
			}
		},
		"constType": func(s string) string {
			cs, ok := jsConstTypes[s]
			if ok {
				return cs
			} else {
				return "string"
			}
		},
	})
}

// getOutputFileName generates an output filename from the babel file path.
func (gen *jsGenerator) getOutputFileName(pidl *idl.Idl, repl string, file string) string {
	s := strings.TrimSuffix(filepath.Base(pidl.Filename), ".babel")
	if !strings.HasSuffix(strings.ToLower(s), strings.ToLower(repl)) {
		s += repl
	}
	//namespaceDirectories := gen.args.Options["output"]
	pkgPath := ""
	pkgPath = filepath.Join(strings.Split(pidl.Namespaces["js"], ".")...)
	/*
		if (namespaceDirectories != "") {
			if (namespaceDirectories == "ns-flat") {
				pkgPath = pidl.Namespaces["js"]
			} else if (namespaceDirectories == "ns-nested") {

			}
		}
	*/
	return filepath.Join(gen.args.OutputDir, pkgPath, file)
}

// GenerateCode generates the source code for the given IDL. It returns an array
// of the generated file names and an error indicator.
func (gen *jsGenerator) GenerateCode(pidl *idl.Idl) ([]string, error) {
	outFnames := make([]string, 0)

	if len(pidl.Consts) > 0 || len(pidl.Enums) > 0 || len(pidl.Structs) > 0 {
		if gen.args.GenModel {
			modelFnames, err := gen.GenCodeInternal(pidl, gen.getOutputFileName(pidl, "model", "model.js"), "model.js")
			if err != nil {
				return nil, err
			}
			outFnames = append(outFnames, modelFnames...)

			pckgFnames, err := gen.GenCodeInternal(pidl, gen.getOutputFileName(pidl, "model", "package.json"), "package.json")
			if err != nil {
				return nil, err
			}
			outFnames = append(outFnames, pckgFnames...)
		}
	}

	//We don't need to generate the service interfaces and the client classes in there is no services defiend in the current file
	if len(pidl.Services) > 0 {
		if gen.args.GenServer {
			serverFnames, err := gen.GenCodeInternal(pidl, gen.getOutputFileName(pidl, "service", "service.js"), "service.js")
			if err != nil {
				return nil, err
			}
			outFnames = append(outFnames, serverFnames...)
			serverImplFnames, err := gen.GenCodeInternal(pidl, gen.getOutputFileName(pidl, "service", "service-impl.js"), "service-impl.js")
			if err != nil {
				return nil, err
			}
			outFnames = append(outFnames, serverImplFnames...)
		}
		if gen.args.GenClient {
			proxyFnames, err := gen.GenCodeInternal(pidl, gen.getOutputFileName(pidl, "client", "client.js"), "client.js")
			if err != nil {
				return nil, err
			}
			outFnames = append(outFnames, proxyFnames...)
		}

	}
	return outFnames, nil
}

func (gen *jsGenerator) GenCodeInternal(pidl *idl.Idl, outFname string, template string) ([]string, error) {
	gen.resetTemplate(pidl)

	if len(pidl.Consts) > 0 || len(pidl.Enums) > 0 || len(pidl.Structs) > 0 || len(pidl.Services) > 0 {

		err := os.MkdirAll(filepath.Dir(outFname), os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("can't create output dir: %w", err)
		}
		os.Remove(outFname)
		outFile, err := os.Create(outFname)
		if err != nil {
			return nil, fmt.Errorf("can't create output file: %w", err)
		}
		defer outFile.Close()
		err = gen.templates.ExecuteTemplate(outFile, template, pidl)
		if err != nil {
			return nil, fmt.Errorf("error executing template: %w", err)
		}
		return []string{outFname}, nil
	} else {
		return []string{}, nil
	}
}
