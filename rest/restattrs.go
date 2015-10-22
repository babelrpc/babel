package rest

import (
	"errors"
	"github.com/babelrpc/babel/idl"
	"strconv"
	"strings"
)

// Error definitions
var (
	ErrMutlipleOps             = errors.New("Only one Op attribute is allowed per method.")
	ErrMutlipleParms           = errors.New("Only one Parm attribute is allowed per method parameter.")
	ErrParmInDataType          = errors.New("Parm.In should be a string with value query, header, path, formData, or body.")
	ErrParmRequiredDataType    = errors.New("Parm.Required should be a bool.")
	ErrParmFormatDataType      = errors.New("Parm.Format should be a string with value csv, ssv, tsv, pipes, or multi.")
	ErrParmNameDataType        = errors.New("Parm.Name should be a string.")
	ErrOpDeprecatedDataType    = errors.New("Op.Deprecated should be a bool.")
	ErrOpHideDataType          = errors.New("Op.Hide should be a bool.")
	ErrOpPathDataType          = errors.New("Op.Path should be a string.")
	ErrOpMethodDataType        = errors.New("Op.Method should be a string with value GET, PUT, POST, DELETE, OPTIONS, HEAD, or PATCH.")
	ErrResponseCodeDataType    = errors.New("Response.Code should be an int.")
	ErrResponseTypeDataType    = errors.New("Response.Type should be a string with a valid Babel type.")
	ErrResponseDescDataType    = errors.New("Response.Desc should be a string.")
	ErrResponseHeadersDataType = errors.New("Response.Headers should be a string containing a comma-separated list of header names.")
	ErrHeaderTypeDataType      = errors.New("Header.Type should be a string with a valid Babel type.")
	ErrHeaderDescDataType      = errors.New("Header.Desc should be a string.")
	ErrHeaderFormatDataType    = errors.New("Header.Format should be a string with value csv, ssv, tsv, pipes, or multi.")
	ErrHeaderNameDataType      = errors.New("Header.Name should be a string.")
	ErrHeaderNameRequired      = errors.New("Header.Name is required.")
	ErrHeaderTypeRequired      = errors.New("Header.Type is required.")
	ErrResponseCodeRequired    = errors.New("Response.Code is required.")
)

// ReadOp processes the REST attributes on a method
func ReadOp(mth *idl.Method) (*Operation, error) {
	found := false
	op := new(Operation)
	op.Responses = make(map[int]Response)
	op.Path = "/"
	headers := make(map[string]Header)
	for _, a := range mth.Attributes {
		if a.Scope == "rest" {
			switch a.Name {
			case "Op":
				if found {
					return nil, ErrMutlipleOps
				} else {
					found = true
				}
				for _, p := range a.Parameters {
					switch p.Name {
					case "Path":
						s, ok := p.Value.(string)
						if !ok || p.DataType != "string" {
							return nil, ErrOpPathDataType
						}
						op.Path = s
					case "Method":
						s, ok := p.Value.(string)
						if !ok || p.DataType != "string" {
							return nil, ErrOpMethodDataType
						}
						switch strings.ToUpper(s) {
						case "GET":
							op.Method = GET
						case "PUT":
							op.Method = PUT
						case "POST":
							op.Method = POST
						case "DELETE":
							op.Method = DELETE
						case "OPTIONS":
							op.Method = OPTIONS
						case "HEAD":
							op.Method = HEAD
						case "PATCH":
							op.Method = PATCH
						default:
							return nil, ErrOpMethodDataType
						}
					case "Deprecated":
						b, ok := p.Value.(bool)
						if !ok || p.DataType != "bool" {
							return nil, ErrOpDeprecatedDataType
						}
						op.Deprecated = b
					case "Hide":
						b, ok := p.Value.(bool)
						if !ok || p.DataType != "bool" {
							return nil, ErrOpHideDataType
						}
						op.Hide = b
					default:
						return nil, errors.New(p.Name + " is not a valid Op parameter.")
					}
				}
			case "Response":
				var resp Response
				resp.Headers = make(map[string]Header)
				resp.Type = mth.Returns
				var code int = -1
				for _, p := range a.Parameters {
					switch p.Name {
					case "Code":
						i, ok := p.Value.(int64)
						if !ok || p.DataType != "int" {
							return nil, ErrResponseCodeDataType
						}
						code = int(i)
					case "Type":
						s, ok := p.Value.(string)
						if !ok || p.DataType != "string" {
							return nil, ErrResponseTypeDataType
						}
						if s != "" {
							t, err := ParseType(s)
							if err != nil {
								return nil, err
							}
							resp.Type = t
						}
					case "Desc":
						s, ok := p.Value.(string)
						if !ok || p.DataType != "string" {
							return nil, ErrResponseDescDataType
						}
						resp.Desc = s
					case "Headers":
						s, ok := p.Value.(string)
						if !ok || p.DataType != "string" {
							return nil, ErrResponseHeadersDataType
						}
						resp.headerNames = strings.Split(strings.Replace(s, " ", "", 0), ",")
					default:
						return nil, errors.New(p.Name + " is not a valid Response parameter.")
					}
				}
				if code < 0 {
					return nil, ErrResponseCodeRequired
				}
				if _, ok := op.Responses[code]; ok {
					return nil, errors.New("Response code " + strconv.Itoa(code) + " alread specificied on method " + mth.Name)
				}
				op.Responses[code] = resp
			case "Header":
				var name string
				var hdr Header
				for _, p := range a.Parameters {
					switch p.Name {
					case "Name":
						s, ok := p.Value.(string)
						if !ok || p.DataType != "string" {
							return nil, ErrHeaderNameDataType
						}
						name = s
					case "Type":
						s, ok := p.Value.(string)
						if !ok || p.DataType != "string" {
							return nil, ErrHeaderTypeDataType
						}
						t, err := ParseType(s)
						if err != nil {
							return nil, err
						}
						hdr.Type = t
					case "Desc":
						s, ok := p.Value.(string)
						if !ok || p.DataType != "string" {
							return nil, ErrHeaderDescDataType
						}
						hdr.Desc = s
					case "Format":
						s, ok := p.Value.(string)
						if !ok || p.DataType != "string" {
							return nil, ErrParmFormatDataType
						}
						switch strings.ToUpper(s) {
						case "":
							hdr.Format = NONE
						case "CSV":
							hdr.Format = CSV
						case "SSV":
							hdr.Format = SSV
						case "TSV":
							hdr.Format = TSV
						case "PIPES":
							hdr.Format = PIPES
						case "MULTI":
							hdr.Format = MULTI
						default:
							return nil, ErrHeaderFormatDataType
						}
					default:
						return nil, errors.New(p.Name + " is not a valid Header parameter.")
					}
				}
				if name == "" {
					return nil, ErrHeaderNameRequired
				}
				if hdr.Type == nil {
					return nil, ErrHeaderTypeRequired
				}
				headers[name] = hdr
			case "Parm":
				return nil, errors.New(a.Name + " is not valid on a method definition.")
			}
		}
	}
	// make sure default response exists (happens when there is no Op or Response attribute)
	if len(op.Responses) == 0 {
		op.Responses[200] = Response{headerNames: []string{}, Type: mth.Returns, Headers: make(map[string]Header)}
	}
	// finish setup by wiring headers
	for _, rsp := range op.Responses {
		for _, h := range rsp.headerNames {
			if hdr, ok := headers[h]; ok {
				rsp.Headers[h] = hdr
			} else {
				return nil, errors.New("Header " + h + " not defined.")
			}
		}
	}
	return op, nil
}

// ReadParm processes the REST attributes on a method parameter
func ReadParm(fld *idl.Field) (*Parm, error) {
	var parm *Parm
	for _, a := range fld.Attributes {
		if a.Scope == "rest" {
			switch a.Name {
			case "Parm":
				if parm != nil {
					return nil, ErrMutlipleParms
				}
				parm = new(Parm)
				for _, p := range a.Parameters {
					switch p.Name {
					case "In":
						s, ok := p.Value.(string)
						if !ok || p.DataType != "string" {
							return nil, ErrParmInDataType
						}
						switch strings.ToUpper(s) {
						case "QUERY":
							parm.In = QUERY
						case "HEADER":
							parm.In = HEADER
						case "PATH":
							parm.In = PATH
						case "FORMDATA":
							parm.In = FORMDATA
						case "BODY":
							parm.In = BODY
						default:
							return nil, ErrParmInDataType
						}
					case "Required":
						b, ok := p.Value.(bool)
						if !ok || p.DataType != "bool" {
							return nil, ErrParmRequiredDataType
						}
						parm.Required = b
					case "Format":
						s, ok := p.Value.(string)
						if !ok || p.DataType != "string" {
							return nil, ErrParmFormatDataType
						}
						switch strings.ToUpper(s) {
						case "":
							parm.Format = NONE
						case "CSV":
							parm.Format = CSV
						case "SSV":
							parm.Format = SSV
						case "TSV":
							parm.Format = TSV
						case "PIPES":
							parm.Format = PIPES
						case "MULTI":
							parm.Format = MULTI
						default:
							return nil, ErrParmFormatDataType
						}
					case "Name":
						s, ok := p.Value.(string)
						if !ok || p.DataType != "string" {
							return nil, ErrParmNameDataType
						}
						parm.Name = s
					default:
						return nil, errors.New(p.Name + " is not a valid Parm parameter.")
					}
				}
			case "Op", "Header", "Response":
				return nil, errors.New(a.Name + " is not valid on a method parameter.")
			}
		}
	}
	// if not provided, use defaults
	if parm == nil {
		parm = new(Parm)
	}
	return parm, nil
}
