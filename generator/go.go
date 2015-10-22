package generator

import (
	"fmt"
	"github.com/babelrpc/babel/idl"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	// Mapping of IDL types to C# types.
	goTypes = map[string]string{
		"bool":     "bool",
		"byte":     "byte",
		"int8":     "int8",
		"int16":    "int16",
		"int32":    "int32",
		"int64":    "int64",
		"float32":  "float32",
		"float64":  "float64",
		"string":   "string",
		"datetime": "time.Time",
		"decimal":  "big.Rat",
		"char":     "rune",
		"binary":   "[]byte",
		"list":     "[]%s",
		"map":      "map[%s]%s",
	}
)

// goGenerator is the code generator for C#.
type goGenerator struct {
	templateManager
	args *Arguments
}

// formatType returns the Type as a string in format suitable for C#.
func (gen *goGenerator) formatType(t *idl.Type) string {
	var s string
	ms, ok := goTypes[t.Name]
	if !ok {
		ms = t.Name
	}
	if t.IsList() {
		s = fmt.Sprintf(ms, gen.formatType(t.ValueType))
	} else if t.IsMap() {
		s = fmt.Sprintf(ms, goTypes[t.KeyType.Name], gen.formatType(t.ValueType))
	} else if t.IsBinary() {
		s = ms
	} else if t.IsPrimitive() {
		s = "*" + ms
	} else if t.IsEnum(gen.tplRootIdl) {
		s = "*" + ms
	} else if t.IsVoid() {
		s = ""
	} else {
		s = "*" + ms
	}
	return s
}

// fullTypeName returns the fully qualified type name using namespace syntax for go.
// Note: this probably doesn't work.
func (gen *goGenerator) fullTypeName(t *idl.Type) string {
	var s string
	ms, ok := goTypes[t.Name]
	if !ok {
		ms = t.Name
	}
	if t.Name == "list" {
		s = fmt.Sprintf(ms, gen.fullTypeName(t.ValueType))
	} else if t.Name == "map" {
		s = fmt.Sprintf(ms, goTypes[t.KeyType.Name], gen.fullTypeName(t.ValueType))
	} else if t.IsPrimitive() {
		s = "*" + ms
	} else {
		ns := gen.tplRootIdl.NamespaceOf(ms, "go")
		if ns != "" {
			s = fmt.Sprintf("*%s.%s", ns, ms)
		} else {
			s = ms
		}
	}
	return s
}

// fullNameOf returns the fully qualified name of the object with its
// namespace formatted for go.
func (gen *goGenerator) fullNameOf(name string) string {
	ns := gen.tplRootIdl.NamespaceOf(name, "go")
	if ns != "" {
		return fmt.Sprintf("%s.%s", ns, name)
	}
	return name
}

// formatLiteral formats a literal value for the Go parser.
func (gen *goGenerator) formatLiteral(value interface{}, typeName string) string {
	if typeName == "#ref" {
		s, ok := value.(string)
		if ok {
			s = strings.Replace(s, ".", "", -1)
			return fmt.Sprintf("%s", s)
		} else {
			return fmt.Sprintf("%s", value)
		}
	} else if typeName == "char" {
		return fmt.Sprintf("%q", value)
	} else {
		return fmt.Sprintf("%#v", value)
	}
}

func (gen *goGenerator) isVoid(t *idl.Type) bool {
	return t.IsVoid()
}

// init sets up the generator for use and loads the templates.
func (gen *goGenerator) init(args *Arguments) error {
	if !args.GenClient && !args.GenModel && !args.GenServer {
		return fmt.Errorf("Nothing to do!")
	} else if args.ServerType != "" {
		return fmt.Errorf("-servertype does not apply to Go.")
	} else if len(args.Options) > 0 {
		for k, _ := range args.Options {
			switch k {
			default:
				return fmt.Errorf("The %s option is not applicable to language go.", k)
			}
		}
	}
	gen.args = args
	return gen.loadTempates(args.TemplateDir, "go", template.FuncMap{
		"formatType":  func(t *idl.Type) string { return gen.formatType(t) },
		"fullNameOf":  func(name string) string { return gen.fullNameOf(name) },
		"formatValue": func(p *idl.Pair) string { return gen.formatLiteral(p.Value, p.DataType) },
		"isVoid":      func(t *idl.Type) bool { return gen.isVoid(t) },
		"imports": func() []string {
			pkg := gen.tplRootIdl.Namespaces["go"]
			imports := make([]string, 0)
			for _, i := range gen.tplRootIdl.UniqueNamespaces("go") {
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
		"serializerOptions": func(t *idl.Type) string {
			if t.IsList() || t.IsMap() || t.IsBinary() {
				return ",omitempty"
			} else if t.IsEnum(gen.tplRootIdl) {
				return ",omitempty"
			} else if t.Name == "int64" || t.Name == "decimal" {
				return ",string,omitempty"
			} else {
				return ",omitempty"
			}
		},
		"package": func() string {
			ns := gen.tplRootIdl.Namespaces["go"]
			arr := strings.Split(ns, "/")
			return arr[len(arr)-1]
		},
		"notptr": func(s string) string {
			if s[0] == '*' {
				return s[1:]
			} else {
				return s
			}
		},
		"modelUsesType": func(s string) bool {
			l, _ := gen.tplRootIdl.UniqueTypes()
			for _, i := range l {
				if i == s {
					return true
				}
			}
			return false
		},
		"serviceUsesType": func(s string) bool {
			_, l := gen.tplRootIdl.UniqueTypes()
			for _, i := range l {
				if i == s {
					return true
				}
			}
			return false
		},
	})
}

// replaceSuffix generates an output filename from the babel file path after ensuring the path exists.
func (gen *goGenerator) replaceSuffix(pidl *idl.Idl, repl string) (string, error) {
	s := strings.TrimSuffix(filepath.Base(pidl.Filename), ".babel")
	if !strings.HasSuffix(strings.ToLower(s), strings.ToLower(repl)) {
		s += repl
	}
	d := filepath.Join(gen.args.OutputDir, filepath.FromSlash(pidl.Namespaces["go"]))
	// check error
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return "", err
	} else {
		return filepath.Join(d, s+".go"), nil
	}
}

// GenerateCode generates the source code for the given IDL. It returns an array
// of the generated file names and an error indicator.
func (gen *goGenerator) GenerateCode(pidl *idl.Idl) ([]string, error) {
	outFnames := make([]string, 0)

	if len(pidl.Consts) > 0 || len(pidl.Enums) > 0 || len(pidl.Structs) > 0 {
		if gen.args.GenModel {
			fn, err := gen.replaceSuffix(pidl, "Model")
			if err != nil {
				return nil, err
			}
			modelFnames, err := gen.GenCodeInternal(pidl, fn, "model.go")
			if err != nil {
				return nil, err
			}
			outFnames = append(outFnames, modelFnames...)
		}
	}

	//We don't need to generate the service interfaces and the client classes in there is no services defiend in the current file
	if len(pidl.Services) > 0 {
		if gen.args.GenModel {
			fn, err := gen.replaceSuffix(pidl, "Interface")
			if err != nil {
				return nil, err
			}
			serviceFnames, err := gen.GenCodeInternal(pidl, fn, "service.go")
			if err != nil {
				return nil, err
			}
			outFnames = append(outFnames, serviceFnames...)
		}

		/*
			if gen.args.GenClient {
				fn, err := gen.replaceSuffix(pidl, "Client")
				if err != nil {
					return nil, err
				}
				proxyFnames, err := gen.GenCodeInternal(pidl, fn, "client.go")
				if err != nil {
					return nil, err
				}
				outFnames = append(outFnames, proxyFnames...)
			}
		*/

		if gen.args.GenServer {
			fn, err := gen.replaceSuffix(pidl, "Invoker")
			if err != nil {
				return nil, err
			}
			invokerFnames, err := gen.GenCodeInternal(pidl, fn, "invoker.go")
			if err != nil {
				return nil, err
			}
			outFnames = append(outFnames, invokerFnames...)
		}
	}
	return outFnames, nil
}

func (gen *goGenerator) GenCodeInternal(pidl *idl.Idl, outFname string, template string) ([]string, error) {
	gen.resetTemplate(pidl)

	if len(pidl.Consts) > 0 || len(pidl.Enums) > 0 || len(pidl.Structs) > 0 || len(pidl.Services) > 0 {
		outFile, err := os.Create(outFname)
		if err != nil {
			return nil, fmt.Errorf("Can't create output file: %s", err)
		}
		defer outFile.Close()
		err = gen.templates.ExecuteTemplate(outFile, template, pidl)
		if err != nil {
			return nil, fmt.Errorf("Error executing template: %s", err)
		}
		//err = gen.RunGoFmt(outFname)
		//if err != nil {
		//	return nil, fmt.Errorf("Error running go fmt: %s", err)
		//}
		return []string{outFname}, nil
	} else {
		return []string{}, nil
	}
}

func (gen *goGenerator) RunGoFmt(outFname string) error {
	return exec.Command("go", "fmt", outFname).Run()
}
