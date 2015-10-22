/*
	IDL files describe models and web services in a format usable by the Babel tools.
	The idl package defines types for the parse tree - the representation of parsed IDL.

	For more information, see the README.md file.
*/
package idl

import (
	"fmt"
	"strings"
)

// Pair is a name/value pair that includes an optional format string for Sprintf.
// Pairs are used to represent constants and other values in the IDL.
type Pair struct {
	Name     string
	Value    interface{}
	DataType string
}

// Const is a collection of constant value defintions. A Const block has a name
// and optional documentation comments.
type Const struct {
	Comments []string
	Name     string
	Values   []*Pair
}

// Init initializes the Const for use.
func (c *Const) Init() {
	c.Values = make([]*Pair, 0)
}

// Add appends a constant value definition to the Const block.
func (c *Const) Add(name string, value interface{}, dataType string) error {
	for _, v := range c.Values {
		if strings.ToLower(v.Name) == strings.ToLower(name) {
			return fmt.Errorf("Constant redefined: %s.%s", c.Name, name)
		}
	}
	if dataType == "" {
		dataType = "string"
	}
	c.Values = append(c.Values, &Pair{Name: name, Value: value, DataType: dataType})
	return nil
}

// FindValue returns the item with the given name.
func (c *Const) FindValue(name string) *Pair {
	for _, v := range c.Values {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// Enum defines a group of enumerated values. An Enum block has a name and optional
// documentation comments.
type Enum struct {
	Comments []string
	Name     string
	Values   []*Pair
}

// Init initializes the Enum for use.
func (e *Enum) Init() {
	e.Values = make([]*Pair, 0)
}

// Add appends an enumeration value to the Enum block.
func (e *Enum) Add(name string, value int64) error {
	for _, v := range e.Values {
		if strings.ToLower(v.Name) == strings.ToLower(name) {
			return fmt.Errorf("Enumeration redefined: %s.%s", e.Name, name)
		}
	}
	e.Values = append(e.Values, &Pair{Name: name, Value: value, DataType: "int"})
	return nil
}

// FindValue returns the item with the given name.
func (e *Enum) FindValue(name string) *Pair {
	for _, v := range e.Values {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// Attribute defines extra metadata for definition following it. Attributes are
// written similar to C# attributes but have meaning specific to the output
// code generator. Attributes have a name and a collection of name/value pairs
// representing the Attirbute's parameters.
type Attribute struct {
	Name       string
	Parameters []*Pair
	Scope      string
}

// Field defines a structure field, which has optional docmumentation comments, optional
// attributes, a type, and a name.
type Field struct {
	Comments    []string
	Attributes  []*Attribute
	Type        *Type
	Name        string
	Initializer *Pair
}

// Init initializes the Field for use.
func (f *Field) Init() {
	f.Comments = make([]string, 0)
	f.Attributes = make([]*Attribute, 0)
}

// Required determines whether the field is required to have a value assigned.
func (f *Field) Required() bool {
	return false
}

// optional determines whether the field is not required to have a vakue assigned.
func (f *Field) Optional() bool {
	return !f.Required()
}

// IsCollection returns true if the type of the field is a list or map.
func (f *Field) IsCollection() bool {
	return f.Type.IsList() || f.Type.IsMap()
}

// IsList returns true if the type of the field is a list
func (f *Field) IsList() bool {
	return f.Type.IsList()
}

// IsMap returns true if the type of the field is a map
func (f *Field) IsMap() bool {
	return f.Type.IsMap()
}

// SetInitializer assigns the given initial value and type, returning an error if it
// isn't compatible with the type of the field.
func (f *Field) SetInitializer(i interface{}, t string) error {
	switch t {
	case "int":
		if !f.Type.IsInt() {
			return fmt.Errorf("Invalid initialization of %s %s with %s", f.Type, f.Name, t)
		}
	case "float":
		if !f.Type.IsFloat() {
			return fmt.Errorf("Invalid initialization of %s %s with %s", f.Type, f.Name, t)
		}
	case "bool":
		if !f.Type.IsBool() {
			return fmt.Errorf("Invalid initialization of %s %s with %s", f.Type, f.Name, t)
		}
	case "string":
		if !f.Type.IsString() {
			return fmt.Errorf("Invalid initialization of %s %s with %s", f.Type, f.Name, t)
		}
	case "char":
		if !f.Type.IsChar() {
			return fmt.Errorf("Invalid initialization of %s %s with %s", f.Type, f.Name, t)
		}
	case "#ref":
		// These are checked later
	default:
		return fmt.Errorf("Invalid data type of initializer for %s: %s", f.Name, t)
	}
	f.Initializer = &Pair{Value: i, DataType: t, Name: f.Name}
	return nil
}

// CheckInitializer verifies that the initalizer is appropriate for the type
func (f *Field) CheckInitializer(idl *Idl) error {
	if f.Initializer != nil {
		t := f.Initializer.DataType
		switch t {
		case "int":
			if !f.Type.IsInt() {
				return fmt.Errorf("Invalid initialization of %s %s with %s", f.Type, f.Name, t)
			}
		case "float":
			if !f.Type.IsFloat() {
				return fmt.Errorf("Invalid initialization of %s %s with %s", f.Type, f.Name, t)
			}
		case "bool":
			if !f.Type.IsBool() {
				return fmt.Errorf("Invalid initialization of %s %s with %s", f.Type, f.Name, t)
			}
		case "string":
			if !f.Type.IsString() {
				return fmt.Errorf("Invalid initialization of %s %s with %s", f.Type, f.Name, t)
			}
		case "char":
			if !f.Type.IsChar() {
				return fmt.Errorf("Invalid initialization of %s %s with %s", f.Type, f.Name, t)
			}
		case "#ref":
			vals := strings.Split(f.Initializer.Value.(string), ".")
			if len(vals) != 2 {
				return fmt.Errorf("Invalid reference syntax for %s %s: %s", f.Type, f.Name, f.Initializer.Value)
			}
			enum := idl.FindEnum(vals[0])
			cons := idl.FindConst(vals[0])
			if enum != nil {
				if !f.Type.IsEnum(idl) {
					return fmt.Errorf("The field %s %s cannot be initialized with an enumeration: %s", f.Type, f.Name, enum.Name)
				}
				if f.Type.Name != enum.Name {
					return fmt.Errorf("The field %s %s is being initialized with the wrong enumeration: %s", f.Type, f.Name, enum.Name)
				}
				if enum.FindValue(vals[1]) == nil {
					return fmt.Errorf("The field %s %s is being initialized with a missing enumeration value: %s.%s", f.Type, f.Name, enum.Name, vals[1])
				}
			} else if cons != nil {
				v := cons.FindValue(vals[1])
				if v == nil {
					return fmt.Errorf("The field %s %s is being initialized with a missing const value: %s.%s", f.Type, f.Name, cons.Name, vals[1])
				}
				if (v.DataType == "int" && !f.Type.IsInt()) || (v.DataType == "float" && !f.Type.IsFloat()) || (v.DataType == "bool" && !f.Type.IsBool()) || (v.DataType == "string" && !f.Type.IsString()) || (v.DataType == "char" && !f.Type.IsChar()) {
					return fmt.Errorf("The field %s %s is being initialized with a constant of the wrong type: %s.%s", f.Type, f.Name, cons.Name, vals[1])
				}
			} else {
				return fmt.Errorf("The field %s %s is being initialized something we can't find: %s.%s", f.Type, f.Name, vals[0], vals[1])
			}
		default:
			return fmt.Errorf("Invalid data type of initializer for %s: %s", f.Name, t)
		}
	}
	return nil
}

// Struct defines a collection of fields transmitted as part of a service call.
// Structs have optional documenation comments and attributes. They may also
// extend other Structs.
type Struct struct {
	Comments   []string
	Attributes []*Attribute
	Name       string
	Extends    string
	Fields     []*Field
	Abstract   bool
}

// Init initializes the Struct for use.
func (s *Struct) Init() {
	s.Comments = make([]string, 0)
	s.Attributes = make([]*Attribute, 0)
	s.Fields = make([]*Field, 0)
}

// AddField adds a field with the given data type and names to the Struct.
func (s *Struct) AddField(dataType *Type, name string) (*Field, error) {
	for _, fld := range s.Fields {
		if strings.ToLower(fld.Name) == strings.ToLower(name) {
			return nil, fmt.Errorf("Structure field redefined: %s.%s", s.Name, name)
		}
	}
	f := new(Field)
	f.Init()
	f.Name = name
	f.Type = dataType
	s.Fields = append(s.Fields, f)
	return f, nil
}

// HasRequiredFields return true if the struct has 1 or more fields that are considered required.
func (s *Struct) HasRequiredFields() bool {
	hasRequired := false
	for _, fld := range s.Fields {
		if fld.Required() {
			hasRequired = true
		}
	}
	return hasRequired
}

// RequiredFields gets all required fields for this struct
func (s *Struct) RequiredFields() []*Field {
	flds := make([]*Field, 0)
	for _, fld := range s.Fields {
		if fld.Required() {
			flds = append(flds, fld)
		}
	}
	return flds
}

// BaseClasses returns the base classes of the named class ordered
// from top to bottom.
func (s *Struct) BaseClasses(idl *Idl) ([]*Struct, error) {
	result := make([]*Struct, 0)
	me := s
	baseName := me.Extends
	for baseName != "" {
		me := idl.FindStruct(baseName)
		if me == nil {
			return nil, fmt.Errorf("Parent not found: %s", baseName)
		}
		baseName = me.Extends
		result = append([]*Struct{me}, result...)
	}
	return result, nil
}

// SubClasses returns a list of structs that extend this one, in no paticular order.
// BUG WARNING: This can't look down into IDL files that are not currently
// loaded. It can only see subclasses in the current context.
func (s *Struct) SubClasses(idl *Idl) []*Struct {
	result := make([]*Struct, 0)
	found := make(map[string]bool)
	for _, v := range idl.Structs {
		if v.Extends == s.Name {
			result = append(result, v)
		}
	}
	for _, imp := range idl.UniqueImports() {
		for _, v := range imp.Structs {
			_, ok := found[v.Name]
			if !ok {
				if v.Extends == s.Name {
					result = append(result, v)
				}
			}
		}
	}
	return result
}

// Method defines a service method that can be called to communicate with a service.
// Methods have parameters and optional documentation comments and attributes.
type Method struct {
	Comments   []string
	Attributes []*Attribute
	Returns    *Type
	Name       string
	Parameters []*Field
}

// Init initializes a Method for use.
func (m *Method) Init() {
	m.Comments = make([]string, 0)
	m.Attributes = make([]*Attribute, 0)
	m.Parameters = make([]*Field, 0)
}

// HasParameters returns true if there are parameters for this method.
func (m *Method) HasParameters() bool {
	return len(m.Parameters) != 0
}

// AddParameter adds a parameter with the given data type and name to the Method.
func (m *Method) AddParameter(dataType *Type, name string) (*Field, error) {
	for _, parm := range m.Parameters {
		if strings.ToLower(parm.Name) == strings.ToLower(name) {
			return nil, fmt.Errorf("Parameter redefined in method %s: %s", m.Name, name)
		}
	}
	p := new(Field)
	p.Init()
	p.Name = name
	p.Type = dataType
	m.Parameters = append(m.Parameters, p)
	return p, nil
}

// Service defines a web service interface. Services have optional documentation
// comments and attributes, as well as a collection of Methods.
type Service struct {
	Comments   []string
	Attributes []*Attribute
	Name       string
	Methods    []*Method
}

// Init initializes the Service for use.
func (s *Service) Init() {
	s.Comments = make([]string, 0)
	s.Attributes = make([]*Attribute, 0)
	s.Methods = make([]*Method, 0)
}

// AddMethod appends a method of the given return type and name to the service.
// A pointer to the new Method is returned.
func (s *Service) AddMethod(returnType *Type, name string) (*Method, error) {
	for _, meth := range s.Methods {
		if strings.ToLower(meth.Name) == strings.ToLower(name) {
			return nil, fmt.Errorf("Method redefined: %s.%s", s.Name, name)
		}
	}
	m := new(Method)
	m.Init()
	m.Name = name
	m.Returns = returnType
	s.Methods = append(s.Methods, m)
	return m, nil
}
