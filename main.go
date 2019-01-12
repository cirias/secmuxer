package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

func main() {
	cmd := flag.String("cmd", "extract", "Command to execute, encrypt or extract")
	store := flag.String("store", "", "directory of the store")
	password := flag.String("password", "", "password used to decrypt")
	flag.Parse()

	if *cmd == "encrypt" {
		err := encryptSecret(os.Stdin, os.Stdout, *password)
		if err != nil {
			log.Fatalln(err)
		}
	} else if *cmd == "extract" {
		err := execute(os.Stdin, os.Stdout, *store, *password)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatalf("Invalid command %s, only encrypt or extract is supported.", *cmd)
	}
}

func encryptSecret(in io.Reader, out io.Writer, password string) error {
	plain, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	encrypted, err := encrypt(bytes.TrimSpace(plain), password)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%s", encrypted)
	return nil
}

func execute(in io.Reader, out io.Writer, store, password string) error {
	resolve := func(filename string) string {
		return filepath.Join(store, filename)
	}
	secret := func(filename string) (string, error) {
		encrypted, err := ioutil.ReadFile(resolve(filename))
		if err != nil {
			return "", err
		}
		plain, err := decrypt(string(encrypted), password)
		return string(plain), err
	}
	funcMap := template.FuncMap{
		"secret": secret,
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
