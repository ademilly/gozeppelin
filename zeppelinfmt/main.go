package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ademilly/gozeppelin/zeppelin"
)

func getInput(path string) (*os.File, error) {
	if path == "" {
		return os.Stdin, nil
	}

	return os.Open(path)
}

func readInput(path string) (string, error) {
	file, err := getInput(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("could not read file to byte")
	}

	return string(b), nil
}

func main() {
	notename := flag.String("name", "new note", "name for the new note")
	filepath := flag.String("filepath", "",
		"[optional] textfile to transform into new notebook JSON formatted Zeppelin request body",
	)
	flag.Parse()

	content, err := readInput(*filepath)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	newNote := zeppelin.TextToNewNote(*notename, *notename, content)
	b, err := json.MarshalIndent(newNote, "", "  ")
	if err != nil {
		log.Fatalf("could not encode struct: %v", err)
	}
	fmt.Println(string(b))
}
