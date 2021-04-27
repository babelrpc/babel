/*
	The generator package defines code generators for Babel. IDL files describe models
	and web services. The generator package defines generators for the supported
	languages. This allows each language to have its own approach to generating code.

	For more information, see the README.md file.
*/
package generator

import (
	"fmt"

	"github.com/babelrpc/babel/idl"
)

// Arguments holds the intialization data needed by the generators.
type Arguments struct {
	TemplateDir string
	OutputDir   string
	Scopes      []string
	GenServer   bool
	GenClient   bool
	GenModel    bool
	ServerType  string
	Options     map[string]string
}

// Generator defines the interface for code generators.
type Generator interface {
	init(args *Arguments) error
	GenerateCode(pidl *idl.Idl) ([]string, error)
}

// New creates a new Generator for the given language, initializing it with
// the given arguments.
func New(lang string, args *Arguments) (Generator, error) {
	var gen Generator
	switch lang {
	case "java":
		gen = new(javaGenerator)
	case "csharp":
		gen = new(csharpGenerator)
	case "asp":
		gen = new(aspGenerator)
	case "test":
		gen = new(testGenerator)
	case "js":
		gen = new(jsGenerator)
	case "go":
		gen = new(goGenerator)
	case "python", "ruby", "javascript", "php", "ios":
		return nil, fmt.Errorf("%s is not supported...yet", lang)
	default:
		return nil, fmt.Errorf("language %s is not supported", lang)
	}
	err := gen.init(args)
	if err != nil {
		return nil, err
	}
	return gen, nil
}
