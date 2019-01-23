package examples

import (
	"fmt"
	"os"
	"strings"
)

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		fmt.Printf("Usage: %s\n%s\n\n", os.Args[0], strings.Join(arg, "\n"))
		os.Exit(1)
	}
}
