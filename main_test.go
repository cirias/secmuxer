package main

import (
	"os"
)

func ExampleExecute() {
	execute(os.Stdout, ".", "test/template")
	// Output:
	// password: password: {{ resolve "test/template" | sh "cat %s" }}
}
