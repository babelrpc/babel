package rest

import (
	"strings"
	"testing"
)

var (
	testCases = []string{
		"bool",
		"int64",
		"list<int64>",
		"list<C>",
		"C",
		"map<string,int32>",
		" map < string , list <   string > >  ",
		"map<string , map<string, map<string, list<string> > > >",
	}
	failCases = []string{
		" map < string , list <<   string > >  ",
		"list<>",
		"< > ",
		"< > list",
		"< string >",
		"int64  string",
		"list<int64>, string",
	}
)

func TestParseType(t *testing.T) {
	for _, c := range testCases {
		typ, err := ParseType(c)
		if err != nil {
			t.Error(err)
		} else if strings.Replace(c, " ", "", -1) != typ.String() {
			t.Errorf("Orignal does not match final: %s vs %s", c, typ.String())
		}
	}
}

func TestParseFails(t *testing.T) {
	for _, c := range failCases {
		_, err := ParseType(c)
		if err == nil {
			t.Error("Should not have parsed: " + c)
		}
	}
}
