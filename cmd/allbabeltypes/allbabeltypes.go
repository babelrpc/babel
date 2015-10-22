package main

import (
	"fmt"
	"github.com/babelrpc/babel/idl"
	"os"
	"strings"
)

func main() {

	fmt.Println("/*")
	fmt.Println("   Babel file with many combinations of data type usage for testing")
	fmt.Println("*/")

	fmt.Println(`
namespace company.com/Babel/TestCases`)

	// --------------------

	fmt.Println(`
// Static tests - strange delimiters are on purpose for parser testing

/// Enumeration doc comment
enum Foo {
	OFF = 0,
	ON = 1
	UNKWOWN = 2,  // how can you not know?
}

/// Constants doc comment
const Bar {
	SomeInt = 23235,
	SomeFloat = 23434.444;
	SomeString = "Hello, Babel!" // and hello to you...
	SomeBool = true,
	SomeChar = 'A'
	NewLine = "\n";
	Unicode = '本'
	MoreUnicode = '\U00101234',
}

/// Happy struct
struct Happy {
	/// The A
	string A;
	/// The B
	string B;
}`)

	// --------------------
	// STRUCTURE TESTS
	// --------------------

	enums := 1
	consts := 1
	structs := 1
	services := 0
	fields := 8
	methods := 0

	fmt.Print(`
// Structure tests - base with all types

/// A nice AAA struct
abstract struct AAA {`)

	a := append(idl.IdlTypes, "Happy")
	for _, t := range a {
		fmt.Println("\n\t/// Doc Comment for my", t)
		fmt.Printf("\t%s my%s; // here's my %s\n", t, strings.ToTitle(t), strings.ToTitle(t))
		fields++
	}
	structs++

	fmt.Println("}")

	// --------------------

	fmt.Println(`
// Structure tests - extends with lists

/// A nice CCC struct
struct CCC extends AAA {`)

	for _, t := range a {
		fmt.Println("\n\t/// Doc Comment for list of", t)
		fmt.Printf("\tlist<%s> listOf%s; // here's my list of %s\n", t, strings.ToTitle(t), strings.ToTitle(t))
		fields++
	}
	structs++

	fmt.Println("}")

	// --------------------

	fmt.Println(`
// Structure tests - extends with maps

/// A nice EEE struct
struct EEE extends CCC {`)

	for _, tk := range idl.IdlTypes {
		if tk != "binary" {
			for _, tv := range a {
				fmt.Println("\n\t/// Doc Comment for map of", tk, "to", tv)
				fmt.Printf("\tmap<%s,%s> mapOf%sto%s; // here's my map of %s to %s\n", tk, tv, strings.ToTitle(tk), strings.ToTitle(tv), strings.ToTitle(tk), strings.ToTitle(tv))
				fields++
			}
		}
	}
	structs++

	fmt.Println("}")

	// --------------------
	// SERVICE TESTS
	// --------------------

	a = append(a, "EEE")

	fmt.Println(`
// Service tests - return types

/// A nice AAA service
service AaaService {`)

	for _, t := range a {
		fmt.Println("\n\t/// Doc Comment for return", t)
		fmt.Printf("\t%s Return%s(); // here's my method returning %s\n", t, strings.ToTitle(t), strings.ToTitle(t))
		methods++
	}
	services++

	fmt.Println("}")

	// --------------------

	fmt.Println(`
// Service tests - accept types

/// A nice CCC service
service CccService {`)

	for _, t := range a {
		fmt.Println("\n\t/// Doc Comment for accept", t)
		fmt.Printf("\tvoid Accept%s(\n\t\t/// my parameter\n\t\t%s parameter); // here's my method accepting %s\n", strings.ToTitle(t), t, strings.ToTitle(t))
		methods++
	}
	services++

	fmt.Println("}")

	// --------------------
	// SERVICE TESTS - LISTS
	// --------------------

	fmt.Println(`
// Service tests - return lists

/// A nice EEE service
service EeeService {`)

	for _, t := range a {
		fmt.Println("\n\t/// Doc Comment for return list of ", t)
		fmt.Printf("\tlist<%s> ReturnListOf%s(); // here's my method returning list of %s\n", t, strings.ToTitle(t), strings.ToTitle(t))
		methods++
	}
	services++

	fmt.Println("}")

	// --------------------

	fmt.Println(`
// Service tests - accept lists

/// A nice GGG service
service GggService {`)

	for _, t := range a {
		fmt.Println("\n\t/// Doc Comment for accept list of", t)
		fmt.Printf("\tvoid AcceptListOf%s(\n\t\t/// my list parameter\n\t\tlist<%s> parameter); // here's my method accepting list of %s\n", strings.ToTitle(t), t, strings.ToTitle(t))
		methods++
	}
	services++

	fmt.Println("}")

	// --------------------
	// SERVICE TESTS - MAPS
	// --------------------

	fmt.Println(`
// Service tests - return maps

/// A nice III service
service IiiService {`)

	for _, tk := range idl.IdlTypes {
		if tk != "binary" {
			for _, tv := range a {
				fmt.Println("\n\t/// Doc Comment for return map of", tk, "to", tv)
				fmt.Printf("\tmap<%s,%s> ReturnMapOf%sto%s(); // here's my method returning map of %s to %s\n", tk, tv, strings.ToTitle(tk), strings.ToTitle(tv), strings.ToTitle(tk), strings.ToTitle(tv))
				methods++
			}
		}
	}
	services++

	fmt.Println("}")

	// --------------------

	fmt.Println(`
// Service tests - accept maps

/// A nice L service
service LllService {`)

	for _, tk := range idl.IdlTypes {
		if tk != "binary" {
			for _, tv := range a {
				fmt.Println("\n\t/// Doc Comment for accept map of", tk, "to", tv)
				fmt.Printf("\tvoid AcceptMapOf%sto%s(\n\t\t/// my map parameter\n\t\tmap<%s,%s> parameter); // here's my method accepting map of %s to %s\n", strings.ToTitle(tk), strings.ToTitle(tv), tk, tv, strings.ToTitle(tk), strings.ToTitle(tv))
				methods++
			}
		}
	}
	services++

	fmt.Println("}")

	// --------------------
	// Fun test cases
	// --------------------

	fmt.Println(`
/// A fun structure
struct FunStruct {
	/// X defaults to 2
	int32 X = 2;
	/// Y defaults to Pi
	float32 Y = 3.141592;
	/// Z is orbital
	string Z = "Orbital";
	/// E is some kind of Foo object
	Foo E;
	/// A is a map of strings to list of strings
	map<string,list<string>> A;
	/// B is a list of maps of strings to strings
	list<map<string,string>> B;
}

/// A fun service!
service FunService {
	/// Let's have some fun
	FunStruct HaveFun();
	/// Let's have lots of fun
	list<FunStruct> HaveLotsOfFun();
	/// Wait, who is having fun?
	map<string,list<FunStruct>> WhosHavingFun();
	/// Better save the party until later
	void SaveParty(
		/// A map of strings to list of FunStruct
		map<string,list<FunStruct>> parameter);
}
`)

	structs++
	fields += 2
	services++
	methods += 4

	fmt.Fprintf(os.Stderr, "Generator summary:\n%10d enums\n%10d consts\n%10d structs\n%10d fields\n%10d services\n%10d methods\n",
		enums, consts, structs, fields, services, methods)
}
