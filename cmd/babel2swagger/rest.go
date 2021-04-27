package main

import (
	"errors"
	"html"
	"strconv"
	"strings"

	"github.com/babelrpc/babel/idl"
	"github.com/babelrpc/babel/rest"
	"github.com/babelrpc/swagger2"
)

// RestOperation associates the annotations with the IDL method
type RestOperation struct {
	Annotation *rest.Operation
	IdlMethod  *idl.Method
}

// HttpMethodMap is a map of HTTP methods to rest operations
type HttpMethodMap map[rest.HttpMethod]RestOperation

// Path map is a map of paths to HttpMethodMaps. This groups
// all the IDL methods for the same path in one place.
type PathMap map[string]HttpMethodMap

// addRestService adds Swagger service definitions for RESTful services
func addRestService(swag *swagger2.Swagger, midl *idl.Idl, svc *idl.Service) error {
	svcComments := strings.Join(svc.Comments, "\n")

	// sort methods by common paths
	pathmap := make(PathMap)
	for _, mth := range svc.Methods {
		// Read annotations
		op, err := rest.ReadOp(mth)
		if err != nil {
			return err
		}
		if !op.Hide {
			methmap, ok := pathmap[op.Path]
			if !ok {
				methmap = make(HttpMethodMap)
				pathmap[op.Path] = methmap
			}
			if _, ok := methmap[op.Method]; ok {
				return errors.New("The " + op.Method.String() + " method is defined multiple times for the same path of " + op.Path)
			}
			methmap[op.Method] = RestOperation{Annotation: op, IdlMethod: mth}
		}
	}

	// Loop through the paths
	for path, methmap := range pathmap {
		var p swagger2.PathItem

		for httpmethod, restop := range methmap {
			op := new(swagger2.Operation)
			switch httpmethod {
			case rest.GET:
				p.Get = op
			case rest.PUT:
				p.Put = op
			case rest.POST:
				p.Post = op
			case rest.DELETE:
				p.Delete = op
			case rest.OPTIONS:
				p.Options = op
			case rest.HEAD:
				p.Head = op
			case rest.PATCH:
				p.Patch = op
			}
			op.Tags = []string{svc.Name}
			mthComments := strings.Join(restop.IdlMethod.Comments, "\n")
			theseArgs := make([]string, 0)
			for _, s := range restop.IdlMethod.Parameters {
				// SWAGGER-BUG: Swagger should be HTML escaping this
				theseArgs = append(theseArgs, html.EscapeString(s.Type.String())+" "+s.Name)
			}
			// SWAGGER-BUG: Swagger should be HTML escaping this
			op.Description = html.EscapeString(restop.IdlMethod.Returns.String()) + " " + restop.IdlMethod.Name + "(" + strings.Join(theseArgs, ", ") + ")"
			if mthComments != "" {
				op.Description += "\n\n" + mthComments
			}
			if svcComments != "" {
				op.Description += "\n\n" + svc.Name + ": " + svcComments
			}
			op.OperationId = svc.Name + "_" + restop.IdlMethod.Name
			op.Summary = svc.Name + "." + restop.IdlMethod.Name
			op.Deprecated = restop.Annotation.Deprecated

			// Add parameters
			op.Parameters = make([]swagger2.Parameter, 0)
			for _, fld := range restop.IdlMethod.Parameters {
				parm, err := rest.ReadParm(fld)
				if err != nil {
					return err
				}
				var p *swagger2.Parameter
				if parm.In == rest.BODY {
					p = fieldToBodyParm(midl, fld)
				} else {
					p = fieldToParm(midl, fld)
				}
				if parm.Name != "" {
					p.Name = parm.Name
				}
				p.In = strings.ToLower(parm.In.String())
				if parm.Format != rest.NONE {
					p.Format = strings.ToLower(parm.Format.String())
				}
				if parm.Required {
					p.Required = new(bool)
					*p.Required = parm.Required
				}
				op.Parameters = append(op.Parameters, *p)
			}

			// Add responses
			op.Responses = make(swagger2.Responses)
			for statuscode, resp := range restop.Annotation.Responses {
				strcode := strconv.Itoa(statuscode)
				if statuscode <= 0 {
					strcode = "default"
				}
				desc := resp.Desc
				if desc == "" {
					desc = "Response of type " + html.EscapeString(resp.Type.String())
				}
				theResp := swagger2.Response{
					Description: desc,
					Schema:      returnsToSchema(midl, resp.Type),
					Headers:     make(swagger2.Headers),
				}
				// Headers
				// SWAGGER-BUG: Swagger-ui doesn't show response headers
				for hdrname, hdr := range resp.Headers {
					theHeader := swagger2.Header{
						Description: hdr.Desc,
						ItemsDef:    *typeToItems(midl, hdr.Type),
					}
					if hdr.Format != rest.NONE {
						theHeader.ItemsDef.CollectionFormat = strings.ToLower(hdr.Format.String())
					}
					theResp.Headers[hdrname] = theHeader
				}
				op.Responses[strcode] = theResp
			}
		}
		swag.Paths[path] = p
	}

	return nil
}
