package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ademilly/gozeppelin/zeppelin"
)

var client zeppelin.Client
var hostname string
var port string

func usage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf(`
endpoints:
  - list => list notebooks available on %s
    `, hostname)))
}

func list(w http.ResponseWriter, r *http.Request) {
	parameters := r.URL.Query()

	username := strings.Join(parameters["username"], "")
	password := strings.Join(parameters["password"], "")

	client, err := zeppelin.NewClient(hostname, username, password)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not connect to %s: %v", hostname, err), http.StatusInternalServerError)
		return
	}
	notebooks, err := client.ListNotebooks()
	if err != nil {
		http.Error(w, fmt.Sprintf("could not list notebooks from %s: %v", hostname, err), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(notebooks.Body, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("could not format response from %s: %v", hostname, err), http.StatusInternalServerError)
	}

	w.Write(b)
}

func main() {
	flag.StringVar(&port, "port", "8080", "port number on which to serve")
	flag.StringVar(&hostname, "hostname", "localhost", "zeppelin service hostname")
	flag.Parse()

	add := fmt.Sprintf("0.0.0.0:%s", port)
	srv := http.NewServeMux()

	srv.HandleFunc("/", usage)
	srv.HandleFunc("/list", list)

	log.Printf("serving on %s", add)
	if err := http.ListenAndServe(add, srv); err != nil {
		log.Fatalf("server stopped: %v\n", err)
	}
}
