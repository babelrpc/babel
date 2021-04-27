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
	testTypes = map[string]string{
		"bool":     "bool",
		"byte":     "byte",
		"int8":     "int8",
		"int16":    "int16",
		"int32":    "int32",
		"int64":    "int64",
		"float32":  "float32",
		"float64":  "float64",
		"string":   "string",
		"datetime": "datetime",
		"decimal":  "decimal",
		"char":     "char",
		"binary":   "binary",
		"list":     "list",
		"map":      "map",
	}
)

// testGenerator is the JSON generator for the babel Http/JSON test harness.
type testGenerator struct {
	templateManager
	args *Arguments
}

// formatType returns the Type as a string in format suitable for Test.
func (gen *testGenerator) formatType(t *idl.Type) string {
	ms, ok := testTypes[t.Name]
	if !ok {
		ms = t.Name
	}
	return ms
}

// getTypeKey return the correct JSON key for the specific type, ie type,ref,or enumRef
func (gen *testGenerator) getTypeKey(t *idl.Type) string {
	var s string
	if t.IsEnum(gen.tplRootIdl) {
		s = "enumRef"
	} else if t.IsStruct(gen.tplRootIdl) {
		s = "ref"
	} else {
		s = "type"
	}
	return s
}

// init sets up the generator for use and loads the templates.
func (gen *testGenerator) init(args *Arguments) error {
	if !args.GenClient || !args.GenModel || !args.GenServer {
		return fmt.Errorf("-model, -client, and -server are not applicable to language test")
	} else if args.ServerType != "" {
		return fmt.Errorf("-servertype is not applicable to language test")
	} else if len(args.Options) > 0 {
		return fmt.Errorf("-options is not applicable to language test")
	}
	gen.args = args
	return gen.loadTempates(args.TemplateDir, "test", template.FuncMap{
		"formatType": func(t *idl.Type) string { return gen.formatType(t) },
		"getTypeKey": func(t *idl.Type) string { return gen.getTypeKey(t) },
		"joinComments": func(s []string) string {
			cmts := gen.expandComments(s)
			cmt := strings.Join(cmts, "\\n")
			cmt = strings.Replace(cmt, "\\", "\\\\", -1)
			cmt = strings.Replace(cmt, "\"", "\\\"", -1)
			cmt = strings.Replace(cmt, "\\\\n", "\\n", -1)
			return cmt
		},
		"allStructs": func() []*idl.Struct {
			s := make([]*idl.Struct, 0)
			s = append(s, gen.tplRootIdl.Structs...)
			for _, i := range gen.tplRootIdl.UniqueImports() {
				s = append(s, i.Structs...)
			}
			return s
		},
		"allEnums": func() []*idl.Enum {
			s := make([]*idl.Enum, 0)
			s = append(s, gen.tplRootIdl.Enums...)
			for _, i := range gen.tplRootIdl.UniqueImports() {
				s = append(s, i.Enums...)
			}
			return s
		},
		// "escapeJSON": func(value string) string {
		//     value = strings.Replace(value, "\"", "\"", -1)
		//     return value
		// },
	})
}

// GenerateCode generates the source code for the given IDL. It returns an array
// of the generated file names and an error indicator.
func (gen *testGenerator) GenerateCode(pidl *idl.Idl) ([]string, error) {

	outFnames := make([]string, 0)
	serviceFnames, err := gen.GenService(pidl)
	if err != nil {
		return nil, err
	}
	outFnames = append(outFnames, serviceFnames...)

	return outFnames, nil

}

func (gen *testGenerator) GenService(pidl *idl.Idl) ([]string, error) {

	gen.resetTemplate(pidl)
	outFnames := make([]string, 0)

	for _, s := range pidl.Services {
		fname := filepath.Join(gen.args.OutputDir, s.Name+".json")
		outFnames = append(outFnames, fname)
		outFile, err := os.Create(fname)
		if err != nil {
			return nil, fmt.Errorf("can't create test output file: %w", err)
		}

		defer outFile.Close()
		err = gen.templates.ExecuteTemplate(outFile, "interface.test", s)
		if err != nil {
			return nil, fmt.Errorf("error executing test template: %w", err)
		}
	}

	return outFnames, nil

}
