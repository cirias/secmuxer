package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

var (
	store = flag.String("store", "", "directory of the store")
)

func main() {
	flag.Parse()

	templateFile := flag.Arg(0)

	err := execute(os.Stdout, *store, templateFile)
	if err != nil {
		log.Fatalln(err)
	}
}

func execute(w io.Writer, store, filename string) error {
	resolve := func(filename string) string {
		return filepath.Join(store, filename)
	}
	funcMap := template.FuncMap{
		"resolve": resolve,
		"sh":      sh,
	}

	t, err := template.New(filepath.Base(filename)).Funcs(funcMap).ParseFiles(filename)
	if err != nil {
		log.Fatalln("could not parse template file:", err)
	}

	return errors.Wrap(t.Execute(w, &struct{}{}), "could not execute")
}

func sh(cmdformat string, a ...interface{}) (string, error) {
	shCmd := fmt.Sprintf(cmdformat, a...)
	cmd := exec.Command("sh", "-c", shCmd)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.Output()
	if err != nil {
		return "", errors.Wrapf(err, "could not run sh command: %s", shCmd)
	}

	return string(stdout), nil
}
