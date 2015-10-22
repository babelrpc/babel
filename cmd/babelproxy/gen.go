// for go generate

package main

// github.com/ancientlore/binder is used to package the web files into the executable.
//go:generate binder -package main -o webcontent.go media/*.png swagger/* swagger/*/*
