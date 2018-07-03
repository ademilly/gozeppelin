package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/ademilly/gozeppelin/zeppelin"
)

const (
	envUsername = "GOZEPPELIN_USERNAME"
	envPassword = "GOZEPPELIN_PASSWORD"
)

type credential struct {
	username string
	password string
}

type actionFunc func(*zeppelin.Client) error

func actions() map[string]actionFunc {
	return map[string]actionFunc{
		"list":     list,
		"new-note": newNote,
	}
}

func getKeys(m map[string]actionFunc) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func retrieveCredentialsFromEnv() credential {
	return credential{
		username: os.Getenv(envUsername),
		password: os.Getenv(envPassword),
	}
}

func printAll(r io.Reader) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
}

func list(client *zeppelin.Client) error {
	res, err := client.ListNotebooks()
	if err != nil {
		return err
	}

	log.Println("Logging response body")
	printAll(res.Body)
	defer res.Body.Close()

	return nil
}

func newNote(client *zeppelin.Client) error {
	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		log.Fatalln(errors.New("standard input is not piped in"))
	}

	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("could not read from stdin: %v\n", err)
	}
	var note zeppelin.NewNoteRequestBody
	err = json.Unmarshal(b, &note)
	if err != nil {
		log.Fatalf("could not unmarshal input: %v", err)
	}

	res, err := client.NewNotebook(note)
	if err != nil {
		return err
	}

	log.Println("Logging response body")
	printAll(res.Body)
	defer res.Body.Close()

	return nil
}

func main() {
	action := flag.String("action", "list", `action to perform; supported: [list, new-note]
		- list notebooks
		- new-note creates a new notebook using content from stdin as text for notebook paragraph
	`)
	hostname := flag.String("hostname", "localhost", "url to zeppelin server")
	flag.Parse()

	if _, ok := actions()[*action]; !ok {
		flag.Usage()
		fmt.Printf("Please choose amongst available actions %s", getKeys(actions()))
		os.Exit(0)
	}

	log.Println("Retrieve credentials from environment variables")
	user := retrieveCredentialsFromEnv()

	log.Println("Creating new Zeppelin client")
	client, err := zeppelin.NewClient(*hostname, user.username, user.password)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Performing action %s\n", *action)
	if err := actions()[*action](client); err != nil {
		log.Fatalln(err)
	}
}
