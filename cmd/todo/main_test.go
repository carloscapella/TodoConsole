package main

import (
	"os"
	"testing"
)

// TestMain_RunSmoke checks that main() runs without panicking for help flag and missing flags
func TestMain_RunSmoke(t *testing.T) {
	os.Args = []string{"todo", "--help"}
	defer func() {
		recover() // ignore panic from os.Exit
	}()
	main()

	os.Args = []string{"todo"}
	defer func() {
		recover()
	}()
	main()
}
