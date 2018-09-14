package main

import (
	"os"
)

func ExampleSh() {
	execute(os.Stdout, ".", "", "test/sh")
	// Output:
	// password: password: {{ resolve "test/sh" | sh "cat %s" }}
}

func ExampleSecret() {
	execute(os.Stdout, ".", "s9cr9t", "test/secret")
	// Output:
	// password: password
}
