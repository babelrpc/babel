package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ancientlore/kubismus"
	"github.com/babelrpc/babel/idl"
	"github.com/babelrpc/babel/parser"
	"github.com/babelrpc/babel/rest"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func loadBabelFiles(args []string) (*idl.Idl, error) {
	// initialize map to track processed files
	processedFiles := make(map[string]bool)

	// create base IDL to aggregate into
	var midl idl.Idl
	midl.Init()
	midl.AddDefaultNamespace("babelrpc.github.io", "Foo/Bar")

	// load specified IDL files into it
	for _, infilePat := range args {
		infiles, err := filepath.Glob(infilePat)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Cannot glob files: %s\n", err))
		}
		if len(infiles) == 0 {
			log.Printf("Warning: No files match \"%s\"\n", infilePat)
		}
		for _, infile := range infiles {
			_, ok := processedFiles[infile]
			if ok {
				log.Printf("Already processed %s\n", infile)
			} else {
				processedFiles[infile] = true
				// fmt.Printf("%s:\n", infile)
				bidl, err := parser.ParseIdl(infile, "test")
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Parsing error in %s: %s\n", infile, err))
				}

				midl.Imports = append(midl.Imports, bidl)
			}
		}
	}

	// validate combined babel
	err := midl.Validate("test")
	if err != nil {
		log.Printf("Warning: Combined IDL does not validate: %s\n", err)
	}

	return &midl, nil
}

func allServices(pidl *idl.Idl) []*idl.Service {
	s := make([]*idl.Service, 0)
	s = append(s, pidl.Services...)
	for _, i := range pidl.UniqueImports() {
		s = append(s, i.Services...)
	}
	return s
}

func addHandlers(router *httprouter.Router, midl *idl.Idl) error {
	count := 0
	re := regexp.MustCompile(`\{(\w+)\}`)
	for _, svc := range allServices(midl) {
		for _, mth := range svc.Methods {
			count++
			op, err := rest.ReadOp(mth)
			if err != nil {
				log.Fatal("Cannot process %s.%s: %s", svc.Name, mth.Name, err)
			}
			if !op.Hide {
				routePath := re.ReplaceAllString(path.Join(conf.RestPath, op.Path), ":$1")
				handle, err := makeHandler(midl, svc, mth)
				if err != nil {
					return err
				}
				router.Handle(op.Method.String(), routePath, handle)
				log.Printf("%s %s -> %s", op.Method.String(), path.Join(conf.RestPath, op.Path), path.Join(conf.BabelPath, svc.Name, mth.Name))
			}
			babelPath := path.Join(conf.BabelPath, svc.Name, mth.Name)
			if !strings.HasPrefix(babelPath, "/") {
				babelPath = "/" + babelPath
			}
			handle, err := makeBabelHandler(midl, svc, mth)
			if err != nil {
				return err
			}
			router.Handle("POST", babelPath, handle)
			log.Printf("%s %s -> %s", op.Method.String(), babelPath, babelPath)
		}
	}
	if count == 0 {
		log.Fatal("No services to process")
	}

	return nil
}

func makeHandler(midl *idl.Idl, svc *idl.Service, mth *idl.Method) (httprouter.Handle, error) {

	destPath := path.Join(conf.BabelPath, svc.Name, mth.Name)
	if !strings.HasPrefix(destPath, "/") {
		destPath = "/" + destPath
	}
	destUrl := conf.BabelProto + "://" + conf.BabelAddr + destPath
	handle := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		kubismus.Metric("Requests", 1, 0)
		// make parameter structure
		req := make(map[string]interface{})
		for _, fld := range mth.Parameters {
			annotation, err := rest.ReadParm(fld)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			nm := fld.Name
			if annotation.Name != "" {
				nm = annotation.Name
			}
			switch annotation.In {
			case rest.QUERY:
				val := r.URL.Query()[nm]
				if len(val) > 0 || annotation.Required {
					v, err := toType(val, midl, fld.Type, annotation.Format)
					if err != nil {
						http.Error(w, err.Error(), 500)
						return
					}
					req[fld.Name] = v
				}
			case rest.HEADER:
				val := r.Header[nm]
				if len(val) > 0 || annotation.Required {
					v, err := toType(val, midl, fld.Type, annotation.Format)
					if err != nil {
						http.Error(w, err.Error(), 500)
						return
					}
					req[fld.Name] = v
				}
			case rest.PATH:
				val := p.ByName(nm)
				if val != "" || annotation.Required {
					v, err := toType([]string{val}, midl, fld.Type, annotation.Format)
					if err != nil {
						http.Error(w, err.Error(), 500)
						return
					}
					req[fld.Name] = v
				}
			case rest.FORMDATA:
				// in theory form data is not legal
				// but this is how it could be processed
				err := r.ParseForm()
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				vals := make(map[string]interface{})
				for k, v := range r.PostForm {
					switch len(v) {
					case 0:
						vals[k] = nil
					case 1:
						vals[k] = v[0]
					default:
						vals[k] = v
					}
				}
				req[fld.Name] = vals
			case rest.BODY:
				enc := json.NewDecoder(r.Body)
				var m interface{}
				err := enc.Decode(&m)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				req[fld.Name] = m
			}
		}
		// post to server
		b, err := json.Marshal(req)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		httpreq, err := http.NewRequest("POST", destUrl, bytes.NewReader(b))
		if err != nil {
			log.Fatal(err)
		}
		for k, v := range r.Header {
			for _, vi := range v {
				httpreq.Header.Add(k, vi)
			}
		}
		httpreq.Header.Set("Content-Type", "application/json")
		//httpreq.Header.Set("Accept", "application/json")
		httpreq.ContentLength = int64(len(b))
		httpreq.Close = false
		doHttp(httpreq, w)
	}
	return handle, nil
}

func one(val []string) string {
	if val == nil || len(val) == 0 {
		return ""
	} else {
		return val[0]
	}
}

func toType(val []string, midl *idl.Idl, typ *idl.Type, fmt rest.ListFmt) (interface{}, error) {
	if typ.IsBool() {
		return strconv.ParseBool(one(val))
	} else if typ.IsInt() {
		if typ.Name == "int64" {
			// int64 is quoted
			return one(val), nil
		} else {
			return strconv.ParseInt(one(val), 10, 32)
		}
	} else if typ.IsString() || typ.IsChar() || typ.IsDatetime() || typ.IsBinary() || typ.IsDecimal() {
		// treat as string - these are all quoted
		return one(val), nil
	} else if typ.IsFloat() {
		return strconv.ParseFloat(one(val), 64)
	} else if typ.IsList() && typ.ValueType.IsPrimitive() {
		var sep string
		switch fmt {
		case rest.CSV:
			sep = ","
		case rest.SSV:
			sep = " "
		case rest.TSV:
			sep = "\t"
		case rest.PIPES:
			sep = "|"
		case rest.MULTI:
		default:
			return nil, errors.New("Format must be specified when lists are encoded.")
		}
		vals := make([]interface{}, 0)
		var arr []string
		if sep != "" {
			arr = strings.Split(one(val), sep)
		} else {
			if val == nil {
				arr = []string{}
			} else {
				arr = val
			}
		}
		for _, s := range arr {
			v, err := toType([]string{s}, midl, typ.ValueType, rest.NONE)
			if err != nil {
				return nil, err
			}
			vals = append(vals, v)
		}
		return vals, nil
	} else if typ.IsEnum(midl) {
		// treat as string - quoted
		return one(val), nil
	}

	// all other types are custom
	// in theory maybe maps could be handled?
	return nil, errors.New("Unexpected type: " + typ.String())
}

func makeBabelHandler(midl *idl.Idl, svc *idl.Service, mth *idl.Method) (httprouter.Handle, error) {

	destPath := path.Join(conf.BabelPath, svc.Name, mth.Name)
	if !strings.HasPrefix(destPath, "/") {
		destPath = "/" + destPath
	}
	destUrl := conf.BabelProto + "://" + conf.BabelAddr + destPath
	handle := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		kubismus.Metric("Requests", 1, 0)
		// post to server
		httpreq, err := http.NewRequest("POST", destUrl, r.Body)
		if err != nil {
			log.Fatal(err)
		}
		for k, v := range r.Header {
			for _, vi := range v {
				httpreq.Header.Add(k, vi)
			}
		}
		httpreq.Header.Set("Content-Type", "application/json")
		//httpreq.Header.Set("Accept", "application/json")
		// httpreq.ContentLength = int64(len(b))
		httpreq.Close = false
		doHttp(httpreq, w)
	}
	return handle, nil
}
