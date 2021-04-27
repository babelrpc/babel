package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/babelrpc/babel/idl"
)

// aspGenerator is the code generator for ASP.
type aspGenerator struct {
	templateManager
	args      *Arguments
	extension string
}

// fullNameOf returns the fully qualified name of the object with its
// namespace formatted for ASP.
func (gen *aspGenerator) fullNameOf(name string) string {
	ns := gen.tplRootIdl.NamespaceOf(name, "asp")
	if ns != "" {
		return fmt.Sprintf("%s%s", ns, name)
	}
	return name
}

// internalType returns the Type as a string in a format useful to the ASP
// runtime library. The type names are changed to have their full namespace.
func (gen *aspGenerator) internalType(t *idl.Type) string {
	var s string
	if t.Name == "list" {
		s = fmt.Sprintf("list<%s>", gen.internalType(t.ValueType))
	} else if t.Name == "map" {
		s = fmt.Sprintf("map<%s,%s>", t.KeyType, gen.internalType(t.ValueType))
	} else if t.IsEnum(gen.tplRootIdl) {
		s = "enum"
	} else if !t.IsUserDefined() {
		s = t.Name
	} else {
		s = fmt.Sprintf("%s%s", gen.tplRootIdl.NamespaceOf(t.Name, "asp"), t.Name)
	}
	return s
}

// renamesList produces a list of type renames for list and map entries
func (gen *aspGenerator) renamesList(t *idl.Type) []string {
	s := make([]string, 0)
	if t.IsList() {
		if t.ValueType.Rename == "" {
			s = append(s, t.ValueType.Name)
		} else {
			s = append(s, t.ValueType.Rename)
		}
		s = append(s, gen.renamesList(t.ValueType)...)
	} else if t.IsMap() {
		if t.KeyType.Rename == "" {
			s = append(s, t.KeyType.Name)
		} else {
			s = append(s, t.KeyType.Rename)
		}
		if t.ValueType.Rename == "" {
			s = append(s, t.ValueType.Name)
		} else {
			s = append(s, t.ValueType.Rename)
		}
		s = append(s, gen.renamesList(t.ValueType)...)
	} else {
		s = append(s, "")
	}
	return s
}

// formatLiteral formats a literal value for the ASP parser.
func (gen *aspGenerator) formatLiteral(value interface{}, typeName string, pidl *idl.Idl) string {
	if typeName == "char" {
		if value.(rune) == '\r' {
			return "vbCr"
		} else if value.(rune) == '\n' {
			return "vbLf"
		} else if value.(rune) == '\t' {
			return "vbTab"
		} else {
			return fmt.Sprintf("\"%c\"", value)
		}
	} else if typeName == "string" {
		s := strings.Replace(value.(string), "\t", "\" & vbTab & \"", -1)
		s = strings.Replace(s, "\r", "\" & vbCr & \"", -1)
		s = strings.Replace(s, "\n", "\" & vbLf & \"", -1)
		s = strings.Replace("\""+s+"\"", " & \"\"", "", -1)
		s = strings.Replace(s, "\"\" & ", "", -1)
		return s
	} else if typeName == "#ref" {
		x := strings.Split(value.(string), ".")
		ns := pidl.NamespaceOf(x[0], "asp")
		return ns + strings.Replace(value.(string), ".", "", 1)
	} else {
		return fmt.Sprintf("%#v", value)
	}
}

// filterAttributes returns the attributes for the enabled scopes.
func (gen *aspGenerator) filterAttributes(attrs []*idl.Attribute) []*idl.Attribute {
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

func (gen *aspGenerator) isVoid(t *idl.Type) bool {
	return t.IsVoid()
}

// init sets up the generator for use and loads the templates.
func (gen *aspGenerator) init(args *Arguments) error {
	gen.extension = "asp"
	if !args.GenClient && !args.GenModel {
		return fmt.Errorf("nothing to do")
	} else if args.ServerType != "" {
		return fmt.Errorf("-servertype is not applicable to language asp")
	} else if len(args.Options) > 0 {
		for k, v := range args.Options {
			switch k {
			case "ext":
				switch v {
				case "asp":
				case "vbs":
				default:
					return fmt.Errorf("ext must be \"asp\" or \"vbs\"")
				}
				gen.extension = v
			default:
				return fmt.Errorf("the %s option is not applicable to language asp", k)
			}
		}
	}
	gen.args = args
	return gen.loadTempates(args.TemplateDir, "asp", template.FuncMap{
		"formatType":   func(t *idl.Type) string { return "" },
		"fullNameOf":   func(name string) string { return gen.fullNameOf(name) },
		"formatValue":  func(p *idl.Pair) string { return gen.formatLiteral(p.Value, p.DataType, gen.tplRootIdl) },
		"filterAttrs":  func(attrs []*idl.Attribute) []*idl.Attribute { return gen.filterAttributes(attrs) },
		"isVoid":       func(t *idl.Type) bool { return gen.isVoid(t) },
		"internalType": func(t *idl.Type) string { return gen.internalType(t) },
		"renames": func(t *idl.Type) string {
			if t.IsCollection() {
				return strings.Join(gen.renamesList(t), ",")
			} else {
				return ""
			}
		},
		"isVbs": func() bool { return gen.extension == "vbs" },
		"isAsp": func() bool { return gen.extension == "asp" },
	})
}

// replaceSuffix generates an output filename from the babel file path.
func (gen *aspGenerator) replaceSuffix(filename, repl string) string {
	s := strings.TrimSuffix(filepath.Base(filename), ".babel")
	if !strings.HasSuffix(strings.ToLower(s), strings.ToLower(repl)) {
		s += repl
	}
	return filepath.Join(gen.args.OutputDir, s+"."+gen.extension)
}

// GenerateCode generates the source code for the given IDL. It returns an array
// of the generated file names and an error indicator.
func (gen *aspGenerator) GenerateCode(pidl *idl.Idl) ([]string, error) {
	outFnames := make([]string, 0)

	if len(pidl.Consts) > 0 || len(pidl.Enums) > 0 || len(pidl.Structs) > 0 {
		if gen.args.GenModel {
			modelFnames, err := gen.GenCodeInternal(pidl, gen.replaceSuffix(pidl.Filename, "Model"), "model.asp")
			if err != nil {
				return nil, err
			}
			outFnames = append(outFnames, modelFnames...)
		}
	}

	//We don't need to generate servive interfaces and client classes in there is no services defiend in the current file
	if len(pidl.Services) > 0 {
		if gen.args.GenClient {
			proxyFnames, err := gen.GenCodeInternal(pidl, gen.replaceSuffix(pidl.Filename, "Client"), "client.asp")
			if err != nil {
				return nil, err
			}
			outFnames = append(outFnames, proxyFnames...)
		}
	}
	return outFnames, nil
}

func (gen *aspGenerator) GenCodeInternal(pidl *idl.Idl, outFname string, template string) ([]string, error) {
	gen.resetTemplate(pidl)

	if len(pidl.Consts) > 0 || len(pidl.Enums) > 0 || len(pidl.Structs) > 0 || len(pidl.Services) > 0 {
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
