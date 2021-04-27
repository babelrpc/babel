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
	// Mapping of IDL types to java types.
	javaTypes = map[string]string{
		"bool":     "boolean",
		"byte":     "byte", // java is signed
		"int8":     "byte",
		"int16":    "short",
		"int32":    "int",
		"int64":    "long",
		"float32":  "float",
		"float64":  "double",
		"string":   "String",
		"datetime": "java.util.Date",
		"decimal":  "java.math.BigDecimal",
		"char":     "char",
		"binary":   "byte[]", // java is signed
		"list":     "java.util.List<%s>",
		"map":      "java.util.Map<%s,%s>",
	}

	// Mapping of IDL types to java class types (nullable).
	nullJavaTypes = map[string]string{
		"bool":     "Boolean",
		"byte":     "Byte", // java is signed
		"int8":     "Byte",
		"int16":    "Short",
		"int32":    "Integer",
		"int64":    "Long",
		"float32":  "Float",
		"float64":  "Double",
		"string":   "String",
		"datetime": "java.util.Date",
		"decimal":  "java.math.BigDecimal",
		"char":     "Character",
		"binary":   "byte[]", // java is signed
		"list":     "java.util.List<%s>",
		"map":      "java.util.Map<%s,%s>",
	}
)

// javaGenerator is the code generator for Java.
type javaGenerator struct {
	templateManager
	args *Arguments
}

// formatListInit returns a string initializer for a list Type
func (gen *javaGenerator) formatListInit(t *idl.Type) string {
	var s string
	ms := "java.util.ArrayList<%s>"
	s = fmt.Sprintf(ms, gen.formatType(t.ValueType))
	return s

}

// formatMapInit returns a string initializer for a map Type
func (gen *javaGenerator) formatMapInit(t *idl.Type) string {
	var s string
	ms := "java.util.HashMap<%s,%s>"
	s = fmt.Sprintf(ms, nullJavaTypes[t.KeyType.Name], gen.formatType(t.ValueType))
	return s

}

// formatType returns the Type as a string in format suitable for Java.
func (gen *javaGenerator) formatType(t *idl.Type) string {
	var s string
	ms, ok := nullJavaTypes[t.Name]
	if !ok {
		ms = t.Name
	}
	if t.Name == "list" {
		s = fmt.Sprintf(ms, gen.formatType(t.ValueType))
	} else if t.Name == "map" {
		s = fmt.Sprintf(ms, nullJavaTypes[t.KeyType.Name], gen.formatType(t.ValueType))
	} else {
		s = ms
	}
	return s
}

func (gen *javaGenerator) isVoid(t *idl.Type) bool {
	return t.IsVoid()
}

func (gen *javaGenerator) getParseType(t *idl.Type) string {
	var s string
	ms, ok := nullJavaTypes[t.Name]
	if !ok {
		ms = t.Name
	}
	if t.Name == "list" {
		ms = nullJavaTypes[t.Name]
		s = fmt.Sprintf(ms, gen.formatType(t.ValueType))
	} else if t.Name == "map" {
		ms = nullJavaTypes[t.Name]
		s = fmt.Sprintf(ms, nullJavaTypes[t.KeyType.Name], gen.formatType(t.ValueType))
	} else {
		s = ms
	}
	return s
}

// fullTypeName returns the fully qualified type name using namespace syntax for Java.
// Note: this probably doesn't work.
func (gen *javaGenerator) fullTypeName(t *idl.Type) string {
	var s string
	ms, ok := nullJavaTypes[t.Name]
	if !ok {
		ms = t.Name
	}
	if t.Name == "list" {
		s = fmt.Sprintf(ms, gen.fullTypeName(t.ValueType))
	} else if t.Name == "map" {
		s = fmt.Sprintf(ms, nullJavaTypes[t.KeyType.Name], gen.fullTypeName(t.ValueType))
	} else {
		ns := gen.tplRootIdl.NamespaceOf(ms, "java")
		if ns != "" {
			s = fmt.Sprintf("%s.%s", ns, ms)
		} else {
			s = ms
		}
	}
	return s
}

// fullNameOf returns the fully qualified name of the object with its
// namespace formatted for Java.
func (gen *javaGenerator) fullNameOf(name string) string {
	ns := gen.tplRootIdl.NamespaceOf(name, "java")
	if ns != "" {
		return fmt.Sprintf("%s.%s", ns, name)
	}
	return name
}

// getterName returns the proper name of a getter method
func (gen *javaGenerator) getterName(f *idl.Field) string {
	s := "get" + gen.toPascalCase(f.Name)
	return s
}

// setterName returns the proper name of a setter method
func (gen *javaGenerator) setterName(f *idl.Field) string {
	s := "set" + gen.toPascalCase(f.Name)
	return s
}

// formatLiteral formats a literal value for the Java parser.
func (gen *javaGenerator) formatLiteral(value interface{}, typeName string) string {

	if typeName == "char" {
		return fmt.Sprintf("%q", value)
	} else {
		if typeName == "#ref" {
			return fmt.Sprintf("%s", value)
		} else {
			return fmt.Sprintf("%#v", value)
		}
	}
}

// formatInitializerLiteral formats a initializer literal value for the Java parser.
func (gen *javaGenerator) formatInitializerLiteral(value interface{}, fieldDataType string, initializerDataType string) string {
	if initializerDataType == "#ref" {
		return fmt.Sprintf("%s", value)
	} else if fieldDataType == "char" {
		return fmt.Sprintf("%q", value)
	} else if fieldDataType == "int64" {
		return fmt.Sprintf("%#v", value) + "L"
	} else if fieldDataType == "float32" {
		return fmt.Sprintf("%#v", value) + "f"
	} else if fieldDataType == "float64" {
		return fmt.Sprintf("%#v", value) + "d"
	} else {
		return fmt.Sprintf("%#v", value)
	}
}

// filterAttributes returns the attributes for the enabled scopes.
func (gen *javaGenerator) filterAttributes(attrs []*idl.Attribute) []*idl.Attribute {
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

// init sets up the generator for use and loads the templates.
func (gen *javaGenerator) init(args *Arguments) error {
	if !args.GenClient && !args.GenModel && !args.GenServer {
		return fmt.Errorf("nothing to do")
	} else if args.ServerType != "" {
		return fmt.Errorf("-servertype is not applicable to language java")
	} else if len(args.Options) > 0 {
		return fmt.Errorf("-options is not applicable to language java")
	}
	gen.args = args
	return gen.loadTempates(args.TemplateDir, "java", template.FuncMap{
		"formatType":  func(t *idl.Type) string { return gen.formatType(t) },
		"fullNameOf":  func(name string) string { return gen.fullNameOf(name) },
		"setterName":  func(f *idl.Field) string { return gen.setterName(f) },
		"getterName":  func(f *idl.Field) string { return gen.getterName(f) },
		"formatValue": func(p *idl.Pair) string { return gen.formatLiteral(p.Value, p.DataType) },
		"formatInitializerLiteral": func(f *idl.Field) string {
			return gen.formatInitializerLiteral(f.Initializer.Value, f.Type.Name, f.Initializer.DataType)
		},
		"filterAttrs": func(attrs []*idl.Attribute) []*idl.Attribute { return gen.filterAttributes(attrs) },
		"parseType":   func(t *idl.Type) string { return gen.getParseType(t) },
		"isVoid":      func(t *idl.Type) bool { return gen.isVoid(t) },
		"constType": func(t string) string {
			var s string
			if t == "string" {
				s = "String"
			} else if t == "bool" {
				s = "boolean"
			} else if t == "char" {
				s = "char"
			} else if t == "float" {
				s = "double"
			} else if t == "int" {
				s = "int"
			}
			return s
		},
		"package": func() string { return gen.tplRootIdl.Namespaces["java"] },
		"imports": func() []string {
			pkg := gen.tplRootIdl.Namespaces["java"]
			imports := make([]string, 0)
			for _, i := range gen.tplRootIdl.UniqueNamespaces("java") {
				if i != pkg {
					imports = append(imports, i)
				}
			}
			return imports
		},
		"formatListInit": func(t *idl.Type) string { return gen.formatListInit(t) },
		"formatMapInit":  func(t *idl.Type) string { return gen.formatMapInit(t) },
	})
}

// GenerateCode generates the source code for the given IDL. It returns an array
// of the generated file names and an error indicator.
func (gen *javaGenerator) GenerateCode(pidl *idl.Idl) ([]string, error) {

	outFnames := make([]string, 0)

	if gen.args.GenModel {
		modelFnames, err := gen.GenModel(pidl)
		if err != nil {
			return nil, err
		}
		outFnames = append(outFnames, modelFnames...)
	}

	if gen.args.GenModel {
		enumFnames, err := gen.GenEnum(pidl)
		if err != nil {
			return nil, err
		}
		outFnames = append(outFnames, enumFnames...)
	}

	if gen.args.GenClient || gen.args.GenServer {
		serviceFnames, err := gen.GenService(pidl)
		if err != nil {
			return nil, err
		}
		outFnames = append(outFnames, serviceFnames...)
	}

	if gen.args.GenModel {
		constFnames, err := gen.GenConst(pidl)
		if err != nil {
			return nil, err
		}
		outFnames = append(outFnames, constFnames...)
	}

	return outFnames, nil
}

func (gen *javaGenerator) GenClient(pidl *idl.Idl) ([]string, error) {
	return nil, nil
}

func (gen *javaGenerator) GenService(pidl *idl.Idl) ([]string, error) {

	gen.resetTemplate(pidl)
	outFnames := make([]string, 0)
	pkgPath := gen.BuildPkg(pidl)

	for _, s := range pidl.Services {
		fname := filepath.Join(gen.args.OutputDir, pkgPath, s.Name+".java")
		outFnames = append(outFnames, fname)
		err := os.MkdirAll(filepath.Join(gen.args.OutputDir, pkgPath), os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("can't create service output package dir: %w", err)
		}
		outFile, err := os.Create(fname)
		if err != nil {
			return nil, fmt.Errorf("can't create service output file: %w", err)
		}

		defer outFile.Close()
		err = gen.templates.ExecuteTemplate(outFile, "interface.java", s)
		if err != nil {
			return nil, fmt.Errorf("error executing service template: %w", err)
		}
	}

	return outFnames, nil

}

func (gen *javaGenerator) GenModel(pidl *idl.Idl) ([]string, error) {

	gen.resetTemplate(pidl)
	outFnames := make([]string, 0)
	pkgPath := gen.BuildPkg(pidl)

	for _, s := range pidl.Structs {
		fname := filepath.Join(gen.args.OutputDir, pkgPath, s.Name+".java")
		outFnames = append(outFnames, fname)
		err := os.MkdirAll(filepath.Join(gen.args.OutputDir, pkgPath), os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("can't create model output package dir: %w", err)
		}
		outFile, err := os.Create(fname)
		if err != nil {
			return nil, fmt.Errorf("can't create model output file: %w", err)
		}

		defer outFile.Close()
		err = gen.templates.ExecuteTemplate(outFile, "model.java", s)
		if err != nil {
			return nil, fmt.Errorf("error executing model template: %w", err)
		}
	}

	return outFnames, nil

}

func (gen *javaGenerator) GenEnum(pidl *idl.Idl) ([]string, error) {

	gen.resetTemplate(pidl)
	outFnames := make([]string, 0)
	pkgPath := gen.BuildPkg(pidl)

	for _, e := range pidl.Enums {
		fname := filepath.Join(gen.args.OutputDir, pkgPath, e.Name+".java")
		outFnames = append(outFnames, fname)
		err := os.MkdirAll(filepath.Join(gen.args.OutputDir, pkgPath), os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("can't create enum output package dir: %w", err)
		}
		outFile, err := os.Create(fname)
		if err != nil {
			return nil, fmt.Errorf("can't create enum output file: %w", err)
		}

		defer outFile.Close()
		err = gen.templates.ExecuteTemplate(outFile, "enum.java", e)
		if err != nil {
			return nil, fmt.Errorf("error executing enum template: %w", err)
		}
	}

	return outFnames, nil

}

func (gen *javaGenerator) GenConst(pidl *idl.Idl) ([]string, error) {

	gen.resetTemplate(pidl)
	outFnames := make([]string, 0)
	pkgPath := gen.BuildPkg(pidl)

	for _, c := range pidl.Consts {
		fname := filepath.Join(gen.args.OutputDir, pkgPath, c.Name+".java")
		outFnames = append(outFnames, fname)
		err := os.MkdirAll(filepath.Join(gen.args.OutputDir, pkgPath), os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("can't create const output package dir: %w", err)
		}
		outFile, err := os.Create(fname)
		if err != nil {
			return nil, fmt.Errorf("can't create const output file: %w", err)
		}

		defer outFile.Close()
		err = gen.templates.ExecuteTemplate(outFile, "constants.java", c)
		if err != nil {
			return nil, fmt.Errorf("error executing const template: %w", err)
		}
	}

	return outFnames, nil

}

func (gen *javaGenerator) BuildPkg(pidl *idl.Idl) string {
	return filepath.Join(strings.Split(pidl.Namespaces["java"], ".")...)
}
