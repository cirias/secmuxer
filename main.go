package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

func main() {
	store := flag.String("store", "", "directory of the store")
	password := flag.String("password", "", "password used to decrypt")
	flag.Parse()

	err := execute(os.Stdin, os.Stdout, *store, *password)
	if err != nil {
		log.Fatalln(err)
	}
}

func execute(in io.Reader, out io.Writer, store, password string) error {
	resolve := func(filename string) string {
		return filepath.Join(store, filename)
	}
	secret := func(filename string) (string, error) {
		return sh("openssl aes-256-cbc -a -d -in %s -out - -pass pass:%s", resolve(filename), password)
	}
	funcMap := template.FuncMap{
		"resolve": resolve,
		"sh":      sh,
		"secret":  secret,
	}

	bs, err := ioutil.ReadAll(in)
	if err != nil {
		return errors.Wrap(err, "could not read from input")
	}

	t, err := template.New("template").Funcs(funcMap).Parse(string(bs))
	if err != nil {
		return errors.Wrap(err, "could not parse template file")
	}

	return errors.Wrap(t.Execute(out, nil), "could not execute")
}

func sh(cmdformat string, a ...interface{}) (string, error) {
	shCmd := fmt.Sprintf(cmdformat, a...)
	cmd := exec.Command("sh", "-c", shCmd)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.Output()
	if err != nil {
		return "", errors.Wrapf(err, "could not run sh command: %s", shCmd)
	}

	return strings.TrimSpace(string(stdout)), nil
}
