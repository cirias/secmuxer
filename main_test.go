package main

import (
	"os"
)

func ExampleSh() {
	f, _ := os.Open("test/sh")
	execute(f, os.Stdout, ".", "")
	// Output:
	// password: password: {{ resolve "test/sh" | sh "cat %s" }}
}

func ExampleSecret() {
	f, _ := os.Open("test/secret")
	execute(f, os.Stdout, ".", "s9cr9t")
	// Output:
	// password: password
}
