package rest

import (
	"github.com/babelrpc/babel/idl"
	"github.com/babelrpc/babel/parser"
	"testing"
)

func parseIdl() (*idl.Idl, error) {
	return parser.ParseIdl("testfile.babel", "test")
}

func TestHeaders(t *testing.T) {
	pidl, err := parseIdl()
	if err != nil {
		t.Fatal(err)
	}
	svc := pidl.FindService("A")
	for _, mth := range svc.Methods {
		if mth.Name == "testHeaders" {
			mp, err := ReadOp(mth)
			if err != nil {
				t.Error("A.testHeaders: " + err.Error())
			} else {
				if mp.Method != GET {
					t.Error("Method should be \"get\" for testHeaders")
				}
				if mp.Deprecated != false {
					t.Error("Deprecated should be false for testHeaders")
				}
				if mp.Path != "/test/headers" {
					t.Error("Path should be \"/test/headers\" for testHeaders")
				}
				if len(mp.Responses) != 3 {
					t.Error("Number of responses should be 3 for testHeaders")
				}
				// 201
				resp, ok := mp.Responses[201]
				if !ok {
					t.Error("Response 201 is missing for testHeaders")
				} else {
					if resp.Type.String() != "list<string>" {
						t.Error("Incorrect 201 response type for testHeaders")
					}
					if len(resp.Headers) != 1 {
						t.Error("Wrong number of headers for response 201 in testHeaders")
					}
					hdr, ok := resp.Headers["Foo"]
					if !ok {
						t.Error("Header Foo missing from response 201 in testHeaders")
					} else {
						if hdr.Desc != "Some header" {
							t.Error("Header Foo has wrong description in testHeaders")
						}
						if hdr.Type.String() != "string" {
							t.Error("Header Foo has wrong type in testHeaders")
						}
						if hdr.Format != NONE {
							t.Error("Header Foo as wrong format in testHeader")
						}
					}
				}
				// 0
				resp, ok = mp.Responses[0]
				if !ok {
					t.Error("Response 0 is missing for testHeaders")
				} else {
					if resp.Type.String() != "string" {
						t.Error("Incorrect 0 response type for testHeaders")
					}
					if len(resp.Headers) != 2 {
						t.Error("Wrong number of headers for response 0 in testHeaders")
					}
					hdr, ok := resp.Headers["Foo"]
					if !ok {
						t.Error("Header Foo missing from response 0 in testHeaders")
					} else {
						if hdr.Desc != "Some header" {
							t.Error("Header Foo has wrong description in testHeaders")
						}
						if hdr.Type.String() != "string" {
							t.Error("Header Foo has wrong type in testHeaders")
						}
						if hdr.Format != NONE {
							t.Error("Header Foo as wrong format in testHeader")
						}
					}
					hdr, ok = resp.Headers["Bar"]
					if !ok {
						t.Error("Header Bar missing from response 0 in testHeaders")
					} else {
						if hdr.Desc != "Some other header" {
							t.Error("Header Bar has wrong description in testHeaders")
						}
						if hdr.Type.String() != "list<string>" {
							t.Error("Header Bar has wrong type in testHeaders")
						}
						if hdr.Format != CSV {
							t.Error("Header Bar as wrong format in testHeader")
						}
					}
				}
				// 202
				resp, ok = mp.Responses[202]
				if !ok {
					t.Error("Response 202 is missing for testHeaders")
				} else {
					if resp.Type.String() != "string" {
						t.Error("Incorrect 202 response type for testHeaders")
					}
					if len(resp.Headers) != 1 {
						t.Error("Wrong number of headers for response 202 in testHeaders")
					}
					hdr, ok := resp.Headers["Bar"]
					if !ok {
						t.Error("Header Bar missing from response 202 in testHeaders")
					} else {
						if hdr.Desc != "Some other header" {
							t.Error("Header Bar has wrong description in testHeaders")
						}
						if hdr.Type.String() != "list<string>" {
							t.Error("Header Bar has wrong type in testHeaders")
						}
						if hdr.Format != CSV {
							t.Error("Header Bar as wrong format in testHeader")
						}
					}
				}
			}
		}
	}
}

func TestDefaults(t *testing.T) {
	pidl, err := parseIdl()
	if err != nil {
		t.Fatal(err)
	}
	svc := pidl.FindService("A")
	for _, mth := range svc.Methods {
		if mth.Name == "testDefaults" {
			mp, err := ReadOp(mth)
			if err != nil {
				t.Error("A.testDefaults: " + err.Error())
			} else {
				if mp.Method != GET {
					t.Error("Method should be \"get\" for testDefaults")
				}
				if mp.Deprecated != false {
					t.Error("Deprecated should be false for testDefaults")
				}
				if mp.Path != "/" {
					t.Error("Path should be \"/\" for testDefaults")
				}
				if len(mp.Responses) != 1 {
					t.Error("Number of responses should be 1 for testDefaults")
				}
				resp, ok := mp.Responses[200]
				if !ok {
					t.Error("Response 200 is missing for testDefaults")
				} else {
					if resp.Type.String() != "list<string>" {
						t.Error("Incorrect response type for testDefaults")
					}
					if len(resp.Headers) != 0 {
						t.Error("Should not be headers for testDefaults")
					}
				}
			}
			for _, fld := range mth.Parameters {
				p, err := ReadParm(fld)
				if err != nil {
					t.Error("A.testDefaults." + fld.Name + ": " + err.Error())
				} else {
					switch fld.Name {
					case "a":
						if p.In != QUERY {
							t.Error("In should be \"query\" for paramter a of method testDefaults")
						}
						if p.Required != false {
							t.Error("Required should be false for parameter a of method testDefaults")
						}
						if p.Format != NONE {
							t.Error("Format shoud be \"none\" for paramter a of method testDefaults")
						}
					}
				}
			}
		}
	}
}

func TestParms(t *testing.T) {
	pidl, err := parseIdl()
	if err != nil {
		t.Fatal(err)
	}
	svc := pidl.FindService("A")
	for _, mth := range svc.Methods {
		if mth.Name == "testParms" {
			mp, err := ReadOp(mth)
			if err != nil {
				t.Error("A.testParms: " + err.Error())
			} else {
				if mp.Method != POST {
					t.Error("Method should be \"post\" for testParms")
				}
				if mp.Deprecated != true {
					t.Error("Deprecated should be true for testParms")
				}
				if mp.Path != "/test/{a}/parms" {
					t.Error("Path should be \"/test/{a}/parms\" for tesParms")
				}
				if len(mp.Responses) != 1 {
					t.Error("Number of responses should be 1 for testParms")
				}
				resp, ok := mp.Responses[202]
				if !ok {
					t.Error("Response 202 is missing for testParms")
				} else {
					if resp.Type.String() != "string" {
						t.Error("Incorrect response type for testParms")
					}
					if len(resp.Headers) != 0 {
						t.Error("Should not be headers for testParms")
					}
				}
			}
			for _, fld := range mth.Parameters {
				p, err := ReadParm(fld)
				if err != nil {
					// c and d should fail
					if fld.Name != "c" && fld.Name != "d" {
						t.Error("A.testParms." + fld.Name + ": " + err.Error())
					}
				} else {
					switch fld.Name {
					case "a":
						if p.In != PATH {
							t.Error("In should be \"path\" for paramter a of method testParms")
						}
						if p.Required != true {
							t.Error("Required should be true for parameter a of method testParms")
						}
						if p.Format != NONE {
							t.Error("Format shoud be \"none\" for paramter a of method testParms")
						}
						if p.Name != "A" {
							t.Error("Name shoud be \"A\" for paramter a of method testParms")
						}
					case "b":
						if p.In != QUERY {
							t.Error("In should be \"query\" for paramter b of method testParms")
						}
						if p.Required != false {
							t.Error("Required should be false for parameter b of method testParms")
						}
						if p.Format != PIPES {
							t.Error("Format shoud be \"pipes\" for paramter b of method testParms")
						}
					case "e":
						if p.In != QUERY {
							t.Error("In should be \"query\" for paramter e of method testParms")
						}
						if p.Required != false {
							t.Error("Required should be false for parameter e of method testParms")
						}
						if p.Format != NONE {
							t.Error("Format shoud be \"pipes\" for paramter e of method testParms")
						}
					}
				}
			}
		}
	}
}
