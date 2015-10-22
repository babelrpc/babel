package main

import (
	"github.com/babelrpc/babel/idl"
	"github.com/babelrpc/swagger2"
	"strings"
)

func allStructs(pidl *idl.Idl) []*idl.Struct {
	s := make([]*idl.Struct, 0)
	s = append(s, pidl.Structs...)
	for _, i := range pidl.UniqueImports() {
		s = append(s, i.Structs...)
	}
	return s
}

func allEnums(pidl *idl.Idl) []*idl.Enum {
	s := make([]*idl.Enum, 0)
	s = append(s, pidl.Enums...)
	for _, i := range pidl.UniqueImports() {
		s = append(s, i.Enums...)
	}
	return s
}

func allServices(pidl *idl.Idl) []*idl.Service {
	s := make([]*idl.Service, 0)
	s = append(s, pidl.Services...)
	for _, i := range pidl.UniqueImports() {
		s = append(s, i.Services...)
	}
	return s
}

func typeToItems(pidl *idl.Idl, t *idl.Type) *swagger2.ItemsDef {
	it := new(swagger2.ItemsDef)
	it.Ref = ""
	if t.IsPrimitive() {
		it.Format = t.String()
		if t.IsInt() || t.IsByte() {
			it.Type = "integer"
			it.Format = "int32"
			if t.Name == "int64" {
				if swagInt && restful {
					// Swagger style int64
					it.Format = "int64"
				} else {
					// Babel style int64
					it.Type = "string"  // ??? Babel quotes large integers to avoid precision loss in JavaScript
					it.Format = "int64" // SWAGGER-CLARIFICATION: is format int64 legal with type string?
				}
			}
		} else if t.IsFloat() {
			it.Type = "number"
			it.Format = "float"
			if t.Name == "float64" {
				it.Format = "double"
			}
		} else if t.IsBool() {
			it.Type = "boolean"
			it.Format = ""
		} else if t.IsDatetime() {
			it.Type = "string"
			it.Format = "date-time"
		} else if t.IsDecimal() {
			it.Type = "string"
			it.Format = ""
		} else if t.IsString() || t.IsChar() {
			it.Type = "string"
			it.Format = ""
		}
	} else if t.IsBinary() {
		it.Type = "string"
		it.Format = "byte"
	} else if t.IsMap() {
		it.Type = "object"
		// hmmm....what to do if keytype is not string?
		// SWAGGER-CLARIFICATION: Does swagger require all key types to be strings?
		it.AdditionalProperties = typeToItems(pidl, t.ValueType)
	} else if t.IsList() {
		it.Type = "array"
		it.Format = ""
		it.Items = typeToItems(pidl, t.ValueType)
	} else if t.IsEnum(pidl) {
		// SWAGGER-BUG: Enums cannot be delared in a schema
		it.Type = "string"
		it.Format = ""
		it.Enum = make([]interface{}, 0)
		e := pidl.FindEnum(t.Name)
		if e != nil {
			for _, x := range e.Values {
				it.Enum = append(it.Enum, x.Name)
			}
		}
	} else {
		// user-defined, struct or enum
		it.Ref = "#/definitions/" + t.Name
	}
	return it
}

func fieldToSchema(pidl *idl.Idl, f *idl.Field) *swagger2.Schema {
	sc := new(swagger2.Schema)
	// sc.Title = f.Name
	sc.Description = strings.Join(f.Comments, "\n")
	it := typeToItems(pidl, f.Type)
	sc.Ref = it.Ref
	sc.Type = it.Type
	sc.Format = it.Format
	sc.ItemsDef.Items = it.Items
	sc.Enum = it.Enum
	sc.AdditionalProperties = it.AdditionalProperties
	return sc
}

// used by rest
func fieldToParm(pidl *idl.Idl, fld *idl.Field) *swagger2.Parameter {
	p := new(swagger2.Parameter)
	p.Name = fld.Name
	p.Description = strings.Join(fld.Comments, "\n")
	it := typeToItems(pidl, fld.Type)
	p.Ref = it.Ref
	p.Type = it.Type
	p.Format = it.Format
	p.ItemsDef.Items = it.Items
	p.Enum = it.Enum
	p.AdditionalProperties = it.AdditionalProperties
	return p
}

// used by rest
func fieldToBodyParm(pidl *idl.Idl, fld *idl.Field) *swagger2.Parameter {
	p := new(swagger2.Parameter)
	p.Name = fld.Name
	p.Description = strings.Join(fld.Comments, "\n")
	p.Schema = fieldToSchema(pidl, fld)
	p.Schema.Description = ""
	return p
}

func parmsToSchema(pidl *idl.Idl, m *idl.Method) *swagger2.Schema {
	sc := new(swagger2.Schema)
	// sc.Title
	sc.Description = strings.Join(m.Comments, "\n")
	sc.Properties = make(map[string]swagger2.Schema)
	sc.Type = "object"
	if m.HasParameters() {
		for _, p := range m.Parameters {
			sc.Properties[p.Name] = *fieldToSchema(pidl, p)
		}
	}
	return sc
}

func returnsToSchema(pidl *idl.Idl, t *idl.Type) *swagger2.Schema {
	sc := new(swagger2.Schema)
	// sc.Title = f.Name
	if t.IsVoid() {
		// nil schema means the operation returns no content
		return nil
	} else {
		it := typeToItems(pidl, t)
		sc.Ref = it.Ref
		sc.Type = it.Type
		sc.Format = it.Format
		sc.ItemsDef.Items = it.Items
		sc.Enum = it.Enum
		sc.AdditionalProperties = it.AdditionalProperties
	}
	return sc
}

func structToSchema(pidl *idl.Idl, st *idl.Struct) *swagger2.Schema {
	sc := new(swagger2.Schema)
	// sc.Title
	sc.Description = strings.Join(st.Comments, "\n")
	sc.Properties = make(map[string]swagger2.Schema)
	sc.Type = "object"
	for _, p := range st.Fields {
		sc.Properties[p.Name] = *fieldToSchema(pidl, p)
	}
	if st.Extends != "" {
		sc.AllOf = make([]swagger2.Schema, 0)

		// SWAGGER-BUG: swagger-js and swagger-ui don't support allOf. Alternate implementation below
		// see https://github.com/swagger-api/swagger-js/issues/188
		if !flatten {
			sc.AllOf = append(sc.AllOf, swagger2.Schema{ItemsDef: swagger2.ItemsDef{Ref: "#/definitions/" + st.Extends}})

			// Attach myself as second schema in AllOf
			sc2 := new(swagger2.Schema)
			sc2.Properties = sc.Properties
			sc2.Type = "object"
			sc.Type = ""
			sc.Properties = nil
			sc.AllOf = append(sc.AllOf, *sc2)
		} else {
			// Alternate implementation due to swagger bug
			bases, err := st.BaseClasses(pidl)
			if err != nil {
				panic("cannot get base classes")
			}
			// add base properties
			for _, b := range bases {
				for _, p := range b.Fields {
					sc.Properties[p.Name] = *fieldToSchema(pidl, p)
				}
			}
		}
	}

	return sc
}

func enumToSchema(pidl *idl.Idl, e *idl.Enum) *swagger2.Schema {
	sc := new(swagger2.Schema)
	// sc.Title = f.Name
	sc.Description = strings.Join(e.Comments, "\n")
	sc.Ref = ""
	sc.Type = "string"
	sc.Format = ""
	sc.Enum = make([]interface{}, 0)
	for _, x := range e.Values {
		sc.Enum = append(sc.Enum, x.Name)
	}
	return sc
}
