package idl

import (
	"fmt"
)

var (
	IdlTypes = []string{
		"bool",     // boolean (true/false)
		"byte",     // unsigned 8-bit integer  (33, 22)
		"int8",     // signed 8-bit integer    (-33, 22)
		"int16",    // signed 16-bit integer   (-9655, 22)
		"int32",    // signed 32-bit integer   (-944332, 22)
		"int64",    // signed 64-bit integer   (-3443434343, 22)
		"float32",  // 32-bit floating point   (3.14)
		"float64",  // 64-bit floating point   (3.141592)
		"string",   // string of utf8 text     ("\tSome text\r\n")
		"datetime", // date-time object        ("2013-09-09T13:44.22.341-05:00")
		"decimal",  // large decimal value     ("300.2415443")
		"char",     // single utf8 character   ("\r", "A")
		"binary",   // byte array or blob      ("YXNhZGFzZAo=")
	}
	IdlContainers = []string{
		"list",
		"map",
	}
)

// Error defines fields that will be logged in a uniform format for build tools to process
// Source (Line, Column): Category (error|warning) Code: Message
type Error struct {
	Source    string // source of the error - program name, file name, etc. Defaults to "babel"
	Line      int    // Line number in the file (optional)
	Column    int    // Column number within the line (optional)
	Category  string // Additional cateogrization "Commmand line", "Parsing", etc. (optional)
	IsWarning bool   // True if the message is a warning message
	Code      int    // Specific error number
	Message   error  // Text of the error message or another error
}

// Implement the error interface
func (e *Error) Error() string {
	prefix := ""
	if e.Source != "" {
		prefix += e.Source
	} else {
		prefix += "babel"
	}
	if e.Line > 0 {
		if e.Column > 0 {
			prefix += fmt.Sprintf("(%d,%d)", e.Line, e.Column)
		} else {
			prefix += fmt.Sprintf("(%d)", e.Line)
		}
	}
	prefix += ": "
	if e.Category != "" {
		prefix += e.Category + " "
	}
	if e.IsWarning {
		prefix += fmt.Sprintf("warning %d: ", e.Code)
	} else {
		prefix += fmt.Sprintf("error %d: ", e.Code)
	}
	return prefix + e.Message.Error()
}

// Sample: return &idl.Error{Code = 500, Message = fmt.Errorf(...)}
