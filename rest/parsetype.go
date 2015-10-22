package rest

import (
	"errors"
	"github.com/babelrpc/babel/idl"
	"regexp"
)

var (
	mapRE   = regexp.MustCompile(`^\s*map\s*<(.*)>\s*$`)
	listRE  = regexp.MustCompile(`^\s*list\s*<(.*)>\s*$`)
	basicRE = regexp.MustCompile(`^\s*(\w+)\s*$`)
)

// ParseType processes a type string and returns the IDL type for it
func ParseType(s string) (*idl.Type, error) {
	var err error
	var t idl.Type
	var m []string
	if m = mapRE.FindStringSubmatch(s); m != nil {
		i := findSplit(m[1])
		if i < 0 {
			return nil, errors.New("Unable to break key/value: " + s)
		}
		t.Name = "map"
		t.KeyType, err = ParseType(m[1][:i])
		if err != nil {
			return nil, err
		}
		if !t.KeyType.IsPrimitive() {
			return nil, errors.New("Key type must be a primitive type: " + s)
		}
		t.ValueType, err = ParseType(m[1][i+1:])
		if err != nil {
			return nil, err
		}
	} else if m = listRE.FindStringSubmatch(s); m != nil {
		t.Name = "list"
		t.ValueType, err = ParseType(m[1])
		if err != nil {
			return nil, err
		}
	} else if m = basicRE.FindStringSubmatch(s); m != nil {
		t.Name = m[1]
	} else {
		return nil, errors.New("Uanble to parse type: " + s)
	}
	return &t, nil
}

func findSplit(s string) int {
	count := 0
	for i, c := range s {
		switch c {
		case '<':
			count++
		case '>':
			count--
		case ',':
			if count == 0 {
				return i
			}
		}
	}
	return -1
}
