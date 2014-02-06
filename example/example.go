// example - shows the basic usage of mango
//
// Description:
//
// It doesn't take much to create manual pages with mango. Just write down
// stuff you want to include in the manual page in a comment at the top
// of your source file, like this. Feel free to add as many sections as you
// want.
package main

import (
	"flag"
)

var (
	optFoo = flag.Bool("foo", false, "this text should show up in the manual page")

	// If the flag definition follows a comment like this, mango uses the
	// comment as description in the manual page.
	optBar = flag.Bool("bar", false, "the above comment should show up in the manual page")
	optBaz = ""
)

func init() {
	// These two calls reference the same variable and will appear
	// grouped in the manual page. Since these aren't boolean flags, mango
	// prints the argument type as well.
	flag.StringVar(&optBaz, "baz", "", "two calls, but one entry in the manual")
	flag.StringVar(&optBaz, "b", "", "two calls, but one entry in the manual")
}

func main() {
	return
}
