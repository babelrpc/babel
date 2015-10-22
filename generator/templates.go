package generator

import (
	"fmt"
	"github.com/babelrpc/babel/idl"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

// LocateTemplateDir determines the default template folder based on the
// location of the currently running executable.
func LocateTemplateDir() string {
	var exePath, err = exec.LookPath(os.Args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", err)
		exePath = os.Args[0]
	}
	s, err := filepath.Abs(exePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", err)
	} else {
		exePath = s
	}
	exePath, _ = filepath.Split(exePath)
	exePath = filepath.ToSlash(exePath)

	if strings.HasPrefix(exePath, "/usr/bin/") {
		exePath = strings.Replace(exePath, "/usr/bin/", "/etc/", 1)
	} else {
		exePath = strings.Replace(exePath, "/bin/", "/etc/", 1)
	}

	return filepath.FromSlash(exePath) + "babeltemplates"
}

// templateManager holds data needed for template processing
type templateManager struct {
	tplRootIdl *idl.Idl
	tplIndent  string
	templates  *template.Template
}

// resetTemplate reinitializes the IDL and indentation settings prior to
// executing templates.
func (gen *templateManager) resetTemplate(pidl *idl.Idl) {
	gen.tplRootIdl = pidl
	gen.tplIndent = ""
}

// toPascalCase converts a string to PascalCase
func (gen *templateManager) toPascalCase(str string) string {

	if len(str) > 0 {
		res := ""
		forceUpper := true
		for _, c := range str {
			if c == '_' {
				forceUpper = true
			} else {
				if forceUpper {
					res += strings.ToUpper(string(c))
					forceUpper = false
				} else {
					res += string(c)
				}
			}
		}
		return res
	}

	return str
}

// toCamelCase converts a string to camelCase
func (gen *templateManager) toCamelCase(str string) string {

	pascalStr := gen.toPascalCase(str)

	r, n := utf8.DecodeRuneInString(pascalStr)
	camelStr := string(unicode.ToLower(r)) + pascalStr[n:]

	return camelStr
}

// ExpandComments comments, in case of multi-line comments a single slice element may have carriage returns
// in that case each line will become a new slice element in the overall comment slice.
func (gen *templateManager) expandComments(s []string) []string {
	cmts := []string{}
	for _, cmtLine := range s {
		tokens := strings.Split(cmtLine, "\n")
		if len(tokens) > 1 {
			for index, cmt := range tokens {
				if len(strings.TrimSpace(cmt)) > 0 || (index != 0 && index < len(tokens)-1) {
					cmts = append(cmts, strings.TrimSpace(cmt))
				}
			}
		} else {
			cmts = append(cmts, strings.TrimSpace(cmtLine))
		}
	}
	return cmts
}

// getFuncMap returns a function map to use in the templates.
func (gen *templateManager) getFuncMap(xtra template.FuncMap) template.FuncMap {
	m := template.FuncMap{
		"indent":       func() string { return gen.tplIndent },
		"setindent":    func(pr string) string { gen.tplIndent = pr; return "" },
		"addindent":    func(pr string) string { gen.tplIndent = gen.tplIndent + pr; return "" },
		"getindent":    func() string { return gen.tplIndent },
		"idl":          func() *idl.Idl { return gen.tplRootIdl },
		"last":         func(x int, a interface{}) bool { return x == reflect.ValueOf(a).Len()-1 },
		"toCamelCase":  func(name string) string { return gen.toCamelCase(name) },
		"toPascalCase": func(name string) string { return gen.toPascalCase(name) },
		"expandComments": func(s []string) []string {
			return gen.expandComments(s)
		},
	}
	if xtra != nil {
		for k, v := range xtra {
			m[k] = v
		}
	}
	return m
}

// loadTemplates loads the templates for the given language that are found
// in the templates directory.
func (gen *templateManager) loadTempates(templatesDir, lang string, xtraFuncs template.FuncMap) error {
	fil, err := os.Open(filepath.Join(templatesDir, lang))
	if err != nil {
		return err
	}
	defer fil.Close()
	info, err := fil.Readdir(256) // assume no more than 256 templates
	tplFiles := make([]string, 0)

	for _, i := range info {
		if !i.IsDir() {
			tplFiles = append(tplFiles, filepath.Join(templatesDir, lang, i.Name()))
		}
	}
	gen.templates = template.Must(template.New("idl").Funcs(gen.getFuncMap(xtraFuncs)).ParseFiles(tplFiles...))
	return nil
}
