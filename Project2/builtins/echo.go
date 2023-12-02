package builtins

import (
	"fmt"
	"strings"
)

// Echo prints the given arguments.
func Echo(args ...string) {
	fmt.Println(strings.Join(args, " "))
}
