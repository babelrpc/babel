package idl

import (
	"fmt"
)

// Type defines an IDL type with a name. Types can be primitive types like strings or
// ints, but can also be Structs or a collection like maps and lists.
//
//   Name         KeyType      ValueType    Description
//   -----------  -----------  -----------  ----------------------------------
//   (primitive)  (empty)      (empty)      Primitive type
//   (user)       (empty)      (empty)      User defined (Struct, Enum)
//   list         (empty)      Type         List of Type
//   map          (primitive)  Type         Map of (primitive) to Type
//   void         (empty)      (empty)      No return type (only for methods)
//
// Note that types can nest.
type Type struct {
	Name      string // map, list, or type name
	KeyType   *Type  // for maps - this will only ever be a basic primitive type
	ValueType *Type  // for maps and lists
	Rename    string // used to rename this type for some serializers
}

// String returns the Type as a string in IDL format.
func (t *Type) String() string {
	var s string
	if t.Name == "list" {
		s = fmt.Sprintf("list<%s>", t.ValueType)
	} else if t.Name == "map" {
		s = fmt.Sprintf("map<%s,%s>", t.KeyType, t.ValueType)
	} else {
		s = fmt.Sprintf("%s", t.Name)
	}
	return s
}

// TagName returns the name of the type for use in various serializers.
func (t *Type) TagName() string {
	var s string
	if t.Rename != "" {
		return t.Rename
	}
	if t.Name == "list" {
		s = fmt.Sprintf("ListOf%s", t.ValueType.TagName())
	} else if t.Name == "map" {
		s = fmt.Sprintf("MapOf%sTo%s", t.KeyType.TagName(), t.ValueType.TagName())
	} else {
		s = t.Name
	}
	return s
}

// IsAbstract checks if any of the struct members are abstract.
func (t *Type) IsAbstract(idl *Idl) bool {
	if t.IsUserDefined() {
		s := idl.FindStruct(t.Name)
		if s == nil {
			return false
		}
		return s.Abstract
	} else if t.IsList() || t.IsMap() {
		return t.ValueType.IsAbstract(idl)
	} else {
		return false
	}
}

// IsPrimitive checks if the Type is one of the primitive types.
func (t *Type) IsPrimitive() bool {
	switch t.Name {
	case "bool", "byte", "int8", "int16", "int32", "int64", "float32", "float64", "string", "datetime", "decimal", "char":
		return true
	default:
		return false
	}
}

// IsBool checks if the Type is a boolean.
func (t *Type) IsBool() bool {
	return t.Name == "bool"
}

// IsDatetime checks if the Type is a datetime.
func (t *Type) IsDatetime() bool {
	return t.Name == "datetime"
}

// IsDecimal checks if the Type is a decimal.
func (t *Type) IsDecimal() bool {
	return t.Name == "decimal"
}

// IsCollection returns true if the Type is a list or map
func (t *Type) IsCollection() bool {
	return t.IsList() || t.IsMap()
}

// IsInt checks if the Type is an integer type.
func (t *Type) IsInt() bool {
	switch t.Name {
	case "byte", "int8", "int16", "int32", "int64":
		return true
	default:
		return false
	}
}

// IsFloat checks if the Type is a floating-point type.
func (t *Type) IsFloat() bool {
	return t.Name == "float32" || t.Name == "float64"
}

// IsString checks if the Type is a string type.
func (t *Type) IsString() bool {
	return t.Name == "string"
}

// IsChar checks if the Type is a character type.
func (t *Type) IsChar() bool {
	return t.Name == "char"
}

// IsList checks if the Type is a list.
func (t *Type) IsList() bool {
	return t.Name == "list"
}

// IsMap checks if the Type is a map.
func (t *Type) IsMap() bool {
	return t.Name == "map"
}

// IsByte checks if the Type is a byte.
func (t *Type) IsByte() bool {
	return t.Name == "byte" || t.Name == "int8"
}

// IsBinary checks if the Type is a byte array.
func (t *Type) IsBinary() bool {
	return t.Name == "binary"
}

// IsEnum checks if the Type is an Enum that has been defined.
func (t *Type) IsEnum(idl *Idl) bool {
	return t.IsUserDefined() && idl.FindEnum(t.Name) != nil
}

// IsStruct checks if the Type is a Struct that has been defined.
func (t *Type) IsStruct(idl *Idl) bool {
	return t.IsUserDefined() && idl.FindStruct(t.Name) != nil
}

// IsUserDefined checks if the Type is a user-defined type.
func (t *Type) IsUserDefined() bool {
	return !t.IsPrimitive() && !t.IsList() && !t.IsMap() && !t.IsBinary() && !t.IsVoid()
}

// IsVoid checks if the Type is a void (used only for returns from functions).
func (t *Type) IsVoid() bool {
	return t.Name == "void"
}

// Check validates that the Type has been defined.
func (t *Type) Check(idl *Idl) error {
	if t.IsUserDefined() {
		if !t.IsStruct(idl) && !t.IsEnum(idl) {
			return fmt.Errorf("Type %s is not defined", t.Name)
		}
	} else if t.IsList() || t.IsMap() {
		return t.ValueType.Check(idl)
	}
	return nil
}
