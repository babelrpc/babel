package idl

import (
	"fmt"
	"path"
	"strings"
)

// Idl represents the parse tree of an IDL document.
type Idl struct {
	Comments   []string
	Filename   string
	Imports    []*Idl
	Namespaces map[string]string
	Consts     []*Const
	Enums      []*Enum
	Structs    []*Struct
	Services   []*Service
}

// Init initializes the Idl for use.
func (idl *Idl) Init() {
	idl.Comments = make([]string, 0)
	idl.Imports = make([]*Idl, 0)
	idl.Namespaces = make(map[string]string)
	idl.Consts = make([]*Const, 0)
	idl.Enums = make([]*Enum, 0)
	idl.Structs = make([]*Struct, 0)
	idl.Services = make([]*Service, 0)
}

// AddImport appends an imported IDL file to this Idl object.
// Imports could be repeated down the subtrees.
func (idl *Idl) AddImport(fpath string) (*Idl, error) {
	cpath := path.Clean(fpath)
	//fmt.Printf("%s -> %s\n", fpath, cpath)
	for _, itm := range idl.Imports {
		if strings.ToLower(itm.Filename) == strings.ToLower(cpath) {
			return nil, fmt.Errorf("double import of \"%s\"", fpath)
		}
	}
	impIdl := new(Idl)
	impIdl.Init()
	impIdl.Filename = cpath
	idl.Imports = append(idl.Imports, impIdl)
	return impIdl, nil
}

// AddNamespace appends a namespace for the given language.
func (idl *Idl) AddNamespace(language, ns string) error {
	_, ok := idl.Namespaces[language]
	if ok {
		_, ok2 := idl.Namespaces["#default"]
		if !ok2 {
			return fmt.Errorf("Namespace redefined: %s", language)
		}
	}
	idl.Namespaces[language] = strings.TrimSpace(ns)
	return nil
}

// AddDefaultNamespace appends a namespace for all languages, customized for each.
func (idl *Idl) AddDefaultNamespace(domain, ns string) error {
	_, ok := idl.Namespaces["#default"]
	if ok {
		return fmt.Errorf("Default namespace already defined")
	}
	idl.Namespaces["#default"] = strings.TrimSpace(domain) + "/" + strings.TrimSpace(ns)

	domarr := strings.Split(domain, ".")
	if len(domarr) < 2 {
		return fmt.Errorf("Default namespace domain seems too short: %s", domain)
	}
	for i, s := range domarr {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			return fmt.Errorf("Bad domain: %s", domain)
		}
		domarr[i] = strings.ToUpper(s[0:1]) + s[1:]
	}
	for i, j := 0, len(domarr)-1; i < j; i, j = i+1, j-1 {
		domarr[i], domarr[j] = domarr[j], domarr[i]
	}

	arr := strings.Split(ns, "/")
	if len(arr) < 1 {
		return fmt.Errorf("Default namespace seems too short: %s", ns)
	}
	for i, s := range arr {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			return fmt.Errorf("Bad namespace: %s", ns)
		}
		arr[i] = strings.ToUpper(s[0:1]) + s[1:]
	}

	// namespace company.com One.Two -> com.company.one.two
	_, ok = idl.Namespaces["java"]
	if !ok {
		idl.Namespaces["java"] = strings.ToLower(strings.Join(domarr, ".") + "." + strings.Join(arr, "."))
	}

	// namespace company.com One.Two -> Company.One.Two
	_, ok = idl.Namespaces["csharp"]
	if !ok {
		idl.Namespaces["csharp"] = domarr[len(domarr)-1] + "." + strings.Join(arr, ".")
	}

	// namespace company.com One.Two -> OneTwo (prefix on classes)
	_, ok = idl.Namespaces["asp"]
	if !ok {
		idl.Namespaces["asp"] = strings.Join(arr, "")
	}

	// The following may need to be tweaked as the language is implemented

	// namespace company.com One.Two -> one.two (prefix on classes)
	_, ok = idl.Namespaces["js"]
	if !ok {
		idl.Namespaces["js"] = strings.ToLower(strings.Join(arr, "."))
	}

	// namespace company.com One.Two -> OneTwo (prefix on classes)
	_, ok = idl.Namespaces["python"]
	if !ok {
		idl.Namespaces["python"] = strings.Join(arr, "")
	}

	// namespace company.com One.Two -> OneTwo (prefix on classes)
	_, ok = idl.Namespaces["ruby"]
	if !ok {
		idl.Namespaces["ruby"] = strings.Join(arr, "")
	}

	// namespace company.com One.Two -> TWO (prefix on classes)
	_, ok = idl.Namespaces["ios"]
	if !ok {
		s := arr[len(arr)-1]
		slen := len(s)
		if slen > 3 {
			slen = 3
		}
		idl.Namespaces["ios"] = strings.ToUpper(s[0:slen])
	}

	// namespace company.com One.Two -> company.com/one/two
	// this one would be manual: company/one/two
	_, ok = idl.Namespaces["go"]
	if !ok {
		idl.Namespaces["go"] = strings.ToLower(strings.Join(domarr, ".") + "/" + strings.Join(arr, "/"))
	}

	// namespace company.com One.Two -> OneTwo (prefix on classes)
	_, ok = idl.Namespaces["php"]
	if !ok {
		idl.Namespaces["php"] = strings.Join(arr, "")
	}

	return nil
}

// AddConst appends a Const definition.
func (idl *Idl) AddConst(name string) (*Const, error) {
	for _, itm := range idl.Consts {
		if strings.ToLower(itm.Name) == strings.ToLower(name) {
			return nil, fmt.Errorf("Constant block redefined: \"%s\"", name)
		}
	}
	c := new(Const)
	c.Init()
	c.Name = name
	idl.Consts = append(idl.Consts, c)
	return c, nil
}

// AddEnum appends an Enum definition.
func (idl *Idl) AddEnum(name string) (*Enum, error) {
	for _, itm := range idl.Enums {
		if strings.ToLower(itm.Name) == strings.ToLower(name) {
			return nil, fmt.Errorf("Enum block redefined: \"%s\"", name)
		}
	}
	e := new(Enum)
	e.Init()
	e.Name = name
	idl.Enums = append(idl.Enums, e)
	return e, nil
}

// AddStruct appends a Struct definition.
func (idl *Idl) AddStruct(name string) (*Struct, error) {
	for _, itm := range idl.Structs {
		if strings.ToLower(itm.Name) == strings.ToLower(name) {
			return nil, fmt.Errorf("Structure redefined: \"%s\"", name)
		}
	}
	s := new(Struct)
	s.Init()
	s.Name = name
	idl.Structs = append(idl.Structs, s)
	return s, nil
}

// AddService appends a Service definition.
func (idl *Idl) AddService(name string) (*Service, error) {
	for _, itm := range idl.Services {
		if strings.ToLower(itm.Name) == strings.ToLower(name) {
			return nil, fmt.Errorf("Service redefined: \"%s\"", name)
		}
	}
	s := new(Service)
	s.Init()
	s.Name = name
	idl.Services = append(idl.Services, s)
	return s, nil
}

// FindConst searches this Idl and imported Idls for the named Const definition.
func (idl *Idl) FindConst(name string) *Const {
	for _, s := range idl.Consts {
		if strings.ToLower(s.Name) == strings.ToLower(name) {
			return s
		}
	}
	for _, i := range idl.Imports {
		s := i.FindConst(name)
		if s != nil {
			return s
		}
	}
	return nil
}

// FindEnum searches this Idl and imported Idls for the named Enum definition.
func (idl *Idl) FindEnum(name string) *Enum {
	for _, s := range idl.Enums {
		if strings.ToLower(s.Name) == strings.ToLower(name) {
			return s
		}
	}
	for _, i := range idl.Imports {
		s := i.FindEnum(name)
		if s != nil {
			return s
		}
	}
	return nil
}

// FindStruct searches this Idl and imported Idls for the named Struct definition.
func (idl *Idl) FindStruct(name string) *Struct {
	for _, s := range idl.Structs {
		if strings.ToLower(s.Name) == strings.ToLower(name) {
			return s
		}
	}
	for _, i := range idl.Imports {
		s := i.FindStruct(name)
		if s != nil {
			return s
		}
	}
	return nil
}

// FindService searches this Idl and imported Idls for the named Service definition.
func (idl *Idl) FindService(name string) *Service {
	for _, s := range idl.Services {
		if strings.ToLower(s.Name) == strings.ToLower(name) {
			return s
		}
	}
	for _, i := range idl.Imports {
		s := i.FindService(name)
		if s != nil {
			return s
		}
	}
	return nil
}

// NamespaceOf finds the named object and returns the namespace
// from the Idl that the object is defined in. Objects may be
// Structs, Enums, Consts, or Services.
func (idl *Idl) NamespaceOf(name, lang string) string {
	for _, x := range idl.Structs {
		if strings.ToLower(x.Name) == strings.ToLower(name) {
			return idl.Namespaces[lang]
		}
	}
	for _, x := range idl.Enums {
		if strings.ToLower(x.Name) == strings.ToLower(name) {
			return idl.Namespaces[lang]
		}
	}
	for _, x := range idl.Consts {
		if strings.ToLower(x.Name) == strings.ToLower(name) {
			return idl.Namespaces[lang]
		}
	}
	for _, x := range idl.Services {
		if strings.ToLower(x.Name) == strings.ToLower(name) {
			return idl.Namespaces[lang]
		}
	}
	for _, i := range idl.Imports {
		m := i.NamespaceOf(name, lang)
		if m != "" {
			return m
		}
	}
	return ""
}

// Validate tests this Idl for collisions, redefinitions, and other problems.
func (idl *Idl) Validate(lang string) error {
	err := idl.checkForCollisions()
	if err != nil {
		return err
	}
	err = idl.checkStructs()
	if err != nil {
		return err
	}
	err = idl.checkTypes()
	if err != nil {
		return err
	}
	err = idl.checkServices()
	if err != nil {
		return err
	}
	err = idl.checkNamespaces(lang)
	if err != nil {
		return err
	}
	return nil
}

// checkNamespaces verfifies that this Idl and all imported Idls specify a
// namespace for the given language.
func (idl *Idl) checkNamespaces(lang string) error {
	if idl.Namespaces[lang] == "" && lang != "test" {
		return fmt.Errorf("Namespaces for %s are required in %s", lang, idl.Filename)
	}
	for _, i := range idl.Imports {
		err := i.checkNamespaces(lang)
		if err != nil {
			return err
		}
	}
	return nil
}

// checkStructs verifies that all structures and their parent classes exist and don't
// override fields.
func (idl *Idl) checkStructs() error {
	for _, s := range idl.Structs {
		tree := make(map[string]bool)
		tree[s.Name] = true
		fields := make(map[string]bool)
		// add first level fields - these are already checked for uniqueness when added
		for _, f := range s.Fields {
			fields[f.Name] = true
			if f.Type.IsAbstract(idl) == true {
				return fmt.Errorf("Field %s.%s uses an abstract type which is not allowed. Polymorphic types are not supported.", s.Name, f.Name)
			}
			if f.Initializer != nil {
				err := f.CheckInitializer(idl)
				if err != nil {
					return err
				}
			}
		}
		// check uniqueness of fields in parent classes
		baseName := s.Extends
		for baseName != "" {
			_, ok := tree[baseName]
			if ok {
				return fmt.Errorf("Inheritance cycle detected: %s", baseName)
			} else {
				tree[baseName] = true
			}
			inner := idl.FindStruct(baseName)
			if inner == nil {
				return fmt.Errorf("Parent not found: %s", baseName)
			}
			for _, fld := range inner.Fields {
				_, ok := fields[fld.Name]
				if ok {
					return fmt.Errorf("Field %s.%s redefined somewhere up to %s", baseName, fld.Name, s.Name)
				}
				fields[fld.Name] = true
			}
			baseName = inner.Extends
		}
	}
	for _, i := range idl.Imports {
		err := i.checkStructs()
		if err != nil {
			return err
		}
	}
	return nil
}

// checkServices verifies that all services don't use parameters improperly.
func (idl *Idl) checkServices() error {
	for _, s := range idl.Services {
		// methods are already checked for uniqueness when added
		for _, m := range s.Methods {
			if m.Returns.IsAbstract(idl) == true {
				return fmt.Errorf("Method %s.%s returns an abstract type which is not allowed. Polymorphic types are not supported.", s.Name, m.Name)
			}
			// parameters are already checked for uniqueness when added
			hasInitializer := false
			for _, p := range m.Parameters {
				if p.Type.IsAbstract(idl) == true {
					return fmt.Errorf("Parameter %s of method %s.%s uses an abstract type which is not allowed. Polymorphic types are not supported.", p.Name, s.Name, m.Name)
				}
				if p.Initializer == nil && hasInitializer {
					return fmt.Errorf("All initialized parameters of method %s.%s must appear at the end of the method. %s is not initialized.", s.Name, m.Name, p.Name)
				}
				if p.Initializer != nil {
					hasInitializer = true
					err := p.CheckInitializer(idl)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	for _, i := range idl.Imports {
		err := i.checkServices()
		if err != nil {
			return err
		}
	}
	return nil
}

// checkTypes verifies that all types are defined in this Idl or its imports.
func (idl *Idl) checkTypes() error {
	for _, s := range idl.Structs {
		for _, f := range s.Fields {
			err := f.Type.Check(idl)
			if err != nil {
				return err
			}
		}
	}
	for _, s := range idl.Services {
		for _, m := range s.Methods {
			for _, p := range m.Parameters {
				err := p.Type.Check(idl)
				if err != nil {
					return err
				}
			}
			err := m.Returns.Check(idl)
			if err != nil {
				return err
			}
		}
	}
	for _, i := range idl.Imports {
		err := i.checkTypes()
		if err != nil {
			return err
		}
	}
	return nil
}

// checkForCollisions checks for redefined items across this Idl and all imports.
func (idl *Idl) checkForCollisions() error {
	m := make(map[string]bool)
	return idl.checkCollisions(m, false)
}

// checkCollisions checks for redefined items across this Idl and all imports.
func (idl *Idl) checkCollisions(data map[string]bool, shallow bool) error {
	for _, itm := range idl.Consts {
		_, ok := data[strings.ToLower(itm.Name)]
		if ok {
			return fmt.Errorf("Constant \"%s\" redefined in \"%s\"", itm.Name, idl.Filename)
		}
		data[strings.ToLower(itm.Name)] = true
	}
	for _, itm := range idl.Enums {
		_, ok := data[strings.ToLower(itm.Name)]
		if ok {
			return fmt.Errorf("Enum \"%s\" redefined in \"%s\"", itm.Name, idl.Filename)
		}
		data[strings.ToLower(itm.Name)] = true
	}
	for _, itm := range idl.Structs {
		_, ok := data[strings.ToLower(itm.Name)]
		if ok {
			return fmt.Errorf("Struct \"%s\" redefined in \"%s\"", itm.Name, idl.Filename)
		}
		data[strings.ToLower(itm.Name)] = true
	}
	for _, itm := range idl.Services {
		_, ok := data[strings.ToLower(itm.Name)]
		if ok {
			return fmt.Errorf("Service \"%s\" redefined in \"%s\"", itm.Name, idl.Filename)
		}
		data[strings.ToLower(itm.Name)] = true
	}
	if !shallow {
		for _, imp := range idl.UniqueImports() {
			err := imp.checkCollisions(data, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// contains returns true if the given array contains the provided string value.
func contains(s []string, v string) bool {
	for _, itm := range s {
		if itm == v {
			return true
		}
	}
	return false
}

// UniqueNamespaces returns a list of all the namespaces used in this Idl and
// all imports for the given language.
func (idl *Idl) UniqueNamespaces(lang string) []string {
	s := make([]string, 0)
	for _, i := range idl.Imports {
		s2 := i.UniqueNamespaces(lang)
		for _, itm := range s2 {
			if !contains(s, itm) && itm != "" {
				s = append(s, itm)
			}
		}
	}
	if !contains(s, idl.Namespaces[lang]) && idl.Namespaces[lang] != "" {
		s = append(s, idl.Namespaces[lang])
	}
	return s
}

// UniqueImports returns the unique list of all imports referenced by this Idl.
func (idl *Idl) UniqueImports() []*Idl {
	// uniqueness of imports within a single IDL already tested
	s := make([]*Idl, 0)
	names := make([]string, 0)
	for _, i := range idl.Imports {
		fname := strings.ToLower(i.Filename)
		if !contains(names, fname) {
			names = append(names, fname)
			s = append(s, i)
			for _, j := range i.UniqueImports() {
				fname = strings.ToLower(j.Filename)
				if !contains(names, fname) {
					names = append(names, fname)
					s = append(s, j)
				}
			}
		}
	}
	return s
}

// UniqueTypes returns a list of all the types in this Idl,
// first used in structs, second used in services
func (idl *Idl) UniqueTypes() ([]string, []string) {
	rs := make([]string, 0)
	rv := make([]string, 0)

	for _, s := range idl.Structs {
		if !contains(rs, s.Name) {
			rs = append(rs, s.Name)
		}
		for _, f := range s.Fields {
			if !contains(rs, f.Type.String()) {
				rs = append(rs, f.Type.String())
			}
		}
	}
	for _, s := range idl.Services {
		for _, m := range s.Methods {
			for _, p := range m.Parameters {
				if !contains(rv, p.Type.String()) {
					rv = append(rv, p.Type.String())
				}
			}
			if !contains(rv, m.Returns.String()) {
				rv = append(rv, m.Returns.String())
			}
		}
	}

	return rs, rv
}
