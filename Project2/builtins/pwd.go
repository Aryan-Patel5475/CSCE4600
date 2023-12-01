package builtins

import (
	"fmt"
	"os"
)

// PrintWorkingDirectory prints the current working directory.
func PrintWorkingDirectory() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(wd)
	return nil
}
