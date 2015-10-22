package generator

import (
	"fmt"
	"github.com/babelrpc/babel/idl"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	// Mapping of IDL types to C# types.
	csharpTypes = map[string]string{
		"bool":     "bool",
		"byte":     "byte",
		"int8":     "sbyte",
		"int16":    "short",
		"int32":    "int",
		"int64":    "long",
		"float32":  "float",
		"float64":  "double",
		"string":   "string",
		"datetime": "DateTime",
		"decimal":  "decimal",
		"char":     "char",
		"binary":   "byte[]",
		"list":     "List<%s>",
		"map":      "Dictionary<%s,%s>",
	}

	// Mapping of constants to their C# type
	constTypes = map[string]string{
		"int":    "int",
		"float":  "double",
		"string": "string",
		"bool":   "bool",
		"char":   "char",
	}
)

// csharpGenerator is the code generator for C#.
type csharpGenerator struct {
	templateManager
	args               *Arguments
	baseControllerName string
}

// formatType returns the Type as a string in format suitable for C#.
func (gen *csharpGenerator) formatType(t *idl.Type) string {
	var s string
	ms, ok := csharpTypes[t.Name]
	if !ok {
		ms = t.Name
	}
	if t.Name == "list" {
		s = fmt.Sprintf(ms, gen.formatType(t.ValueType))
	} else if t.Name == "map" {
		s = fmt.Sprintf(ms, csharpTypes[t.KeyType.Name], gen.formatType(t.ValueType))
	} else if t.IsPrimitive() && t.Name != "string" {
		s = ms + "?"
	} else if t.IsEnum(gen.tplRootIdl) {
		s = ms + "?"
	} else {
		s = ms
	}
	return s
}

// fullTypeName returns the fully qualified type name using namespace syntax for C#.
// Note: this probably doesn't work.
func (gen *csharpGenerator) fullTypeName(t *idl.Type) string {
	var s string
	ms, ok := csharpTypes[t.Name]
	if !ok {
		ms = t.Name
	}
	if t.Name == "list" {
		s = fmt.Sprintf(ms, gen.fullTypeName(t.ValueType))
	} else if t.Name == "map" {
		s = fmt.Sprintf(ms, csharpTypes[t.KeyType.Name], gen.fullTypeName(t.ValueType))
	} else {
		ns := gen.tplRootIdl.NamespaceOf(ms, "csharp")
		if ns != "" {
			s = fmt.Sprintf("%s.%s", ns, ms)
		} else {
			s = ms
		}
	}
	return s
}

// fullNameOf returns the fully qualified name of the object with its
// namespace formatted for C#.
func (gen *csharpGenerator) fullNameOf(name string) string {
	arr := strings.Split(name, ".")
	if len(arr) > 0 {
		ns := gen.tplRootIdl.NamespaceOf(arr[0], "csharp")
		if ns != "" {
			return ns + "." + strings.Join(arr, ".")
		}
	}
	return name
}

// formatLiteral formats a literal value for the C# parser.
func (gen *csharpGenerator) formatLiteral(value interface{}, typeName string) string {
	if typeName == "#ref" {
		return fmt.Sprintf("%s", gen.fullNameOf(value.(string)))
	} else if typeName == "char" {
		return fmt.Sprintf("%q", value)
	} else {
		return fmt.Sprintf("%#v", value)
	}
}

// filterAttributes returns the attributes for the enabled scopes.
func (gen *csharpGenerator) filterAttributes(attrs []*idl.Attribute) []*idl.Attribute {
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

func (gen *csharpGenerator) isVoid(t *idl.Type) bool {
	return t.IsVoid()
}

func (gen *csharpGenerator) isTrivialProperty(t *idl.Type) bool {
	return t.IsPrimitive() || t.Name == "binary" || t.IsEnum(gen.tplRootIdl)
}

// init sets up the generator for use and loads the templates.
func (gen *csharpGenerator) init(args *Arguments) error {
	if !args.GenClient && !args.GenModel && !args.GenServer {
		return fmt.Errorf("Nothing to do!")
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
					return fmt.Errorf("Invalid output option: %s.  Valid options are 'ns-flat' and 'ns-nested'", v)
				}
			default:
				return fmt.Errorf("The %s option is not applicable to language csharp.", k)
			}
		}
	}
	gen.args = args
	return gen.loadTempates(args.TemplateDir, "csharp", template.FuncMap{
		"formatType":        func(t *idl.Type) string { return gen.formatType(t) },
		"fullNameOf":        func(name string) string { return gen.fullNameOf(name) },
		"formatValue":       func(p *idl.Pair) string { return gen.formatLiteral(p.Value, p.DataType) },
		"filterAttrs":       func(attrs []*idl.Attribute) []*idl.Attribute { return gen.filterAttributes(attrs) },
		"isVoid":            func(t *idl.Type) bool { return gen.isVoid(t) },
		"isTrivialProperty": func(t *idl.Type) bool { return gen.isTrivialProperty(t) },
		"usings": func() []string {
			pkg := gen.tplRootIdl.Namespaces["csharp"]
			imports := make([]string, 0)
			for _, i := range gen.tplRootIdl.UniqueNamespaces("csharp") {
				if i != pkg {
					imports = append(imports, i)
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
			cs, ok := constTypes[s]
			if ok {
				return cs
			} else {
				return "string"
			}
		},
	})
}

// getOutputFileName generates an output filename from the babel file path.
func (gen *csharpGenerator) getOutputFileName(pidl *idl.Idl, repl string) string {
	s := strings.TrimSuffix(filepath.Base(pidl.Filename), ".babel")
	if !strings.HasSuffix(strings.ToLower(s), strings.ToLower(repl)) {
		s += repl
	}
	namespaceDirectories := gen.args.Options["output"]
	pkgPath := ""
	if namespaceDirectories != "" {
		if namespaceDirectories == "ns-flat" {
			pkgPath = pidl.Namespaces["csharp"]
		} else if namespaceDirectories == "ns-nested" {
			pkgPath = filepath.Join(strings.Split(pidl.Namespaces["csharp"], ".")...)
		}
	}

	return filepath.Join(gen.args.OutputDir, pkgPath, s+".cs")
}

// GenerateCode generates the source code for the given IDL. It returns an array
// of the generated file names and an error indicator.
func (gen *csharpGenerator) GenerateCode(pidl *idl.Idl) ([]string, error) {
	outFnames := make([]string, 0)

	if len(pidl.Consts) > 0 || len(pidl.Enums) > 0 || len(pidl.Structs) > 0 {
		if gen.args.GenModel {
			modelFnames, err := gen.GenCodeInternal(pidl, gen.getOutputFileName(pidl, "Model"), "model.cs")
			if err != nil {
				return nil, err
			}
			outFnames = append(outFnames, modelFnames...)
		}
	}

	//We don't need to generate the service interfaces and the client classes in there is no services defiend in the current file
	if len(pidl.Services) > 0 {
		if gen.args.GenModel {
			serviceFnames, err := gen.GenCodeInternal(pidl, gen.getOutputFileName(pidl, "Interface"), "service.cs")
			if err != nil {
				return nil, err
			}
			outFnames = append(outFnames, serviceFnames...)
		}

		if gen.args.GenClient {
			proxyFnames, err := gen.GenCodeInternal(pidl, gen.getOutputFileName(pidl, "Client"), "client.cs")
			if err != nil {
				return nil, err
			}
			outFnames = append(outFnames, proxyFnames...)
		}

		if gen.args.GenServer {
			srvType := "mvc"

			if gen.args.ServerType != "" {
				srvType = gen.args.ServerType
			}
			mvcControllerFnames, err := gen.GenCodeInternal(pidl, gen.getOutputFileName(pidl, srvType+"Controller"), srvType+"Controller.cs")
			if err != nil {
				return nil, err
			}
			outFnames = append(outFnames, mvcControllerFnames...)
		}
	}
	return outFnames, nil
}

func (gen *csharpGenerator) GenCodeInternal(pidl *idl.Idl, outFname string, template string) ([]string, error) {
	gen.resetTemplate(pidl)

	if len(pidl.Consts) > 0 || len(pidl.Enums) > 0 || len(pidl.Structs) > 0 || len(pidl.Services) > 0 {

		err := os.MkdirAll(filepath.Dir(outFname), os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("Can't create output dir: %s", err)
		}
		outFile, err := os.Create(outFname)
		if err != nil {
			return nil, fmt.Errorf("Can't create output file: %s", err)
		}
		defer outFile.Close()
		err = gen.templates.ExecuteTemplate(outFile, template, pidl)
		if err != nil {
			return nil, fmt.Errorf("Error executing template: %s", err)
		}
		return []string{outFname}, nil
	} else {
		return []string{}, nil
	}
}
