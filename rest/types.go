package rest

import (
	"github.com/babelrpc/babel/idl"
)

// HTTP methods
type HttpMethod int

const (
	GET HttpMethod = iota
	PUT
	POST
	DELETE
	OPTIONS
	HEAD
	PATCH
)

//go:generate stringer -type=HttpMethod

// List format
type ListFmt int

const (
	NONE  ListFmt = iota // Not applicable
	CSV                  // comma separated values foo,bar
	SSV                  // space separated values foo bar
	TSV                  // tab separated values foo\tbar
	PIPES                // pipe separated values foo|bar
	MULTI                // corresponds to multiple parameter instances instead of multiple values for a single instance foo=bar&foo=baz. This is valid only for parameters in "query" or "formData"
)

//go:generate stringer -type=ListFmt

// Defines where a parameter is used
type ParmIn int

const (
	QUERY    ParmIn = iota // query string
	HEADER                 // header
	PATH                   // in the path
	FORMDATA               // form data
	BODY                   // body
)

//go:generate stringer -type=ParmIn

// Operation defines a RESTful operation
type Operation struct {
	Path       string           // API path
	Method     HttpMethod       // HTTP method
	Deprecated bool             // whether the API is deprecated
	Responses  map[int]Response // Response that can occur
	Hide       bool             // Don't expose this operation from tools
}

// Defines a response that can come back from an API
type Response struct {
	Type        *idl.Type         // Response type
	Headers     map[string]Header // Headers that are returned
	Desc        string            // Response description
	headerNames []string          // value from attributes
}

// Defines an HTTP header
type Header struct {
	Type   *idl.Type // Header type
	Desc   string    // Description
	Format ListFmt   // format of list parameters
}

// Defines a parameter
type Parm struct {
	In       ParmIn  // where the parameter comes from
	Required bool    // whether it is required
	Format   ListFmt // format of list parameters
	Name     string  // used to rename (usually for headers)
}
