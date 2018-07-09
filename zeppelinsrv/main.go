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
  - run  => run notebooks by IDs given as comma separated URL parameter notebookIDs
    `, hostname)))
}

func login(_ http.ResponseWriter, r *http.Request) (*zeppelin.Client, error) {
	parameters := r.URL.Query()

	username := strings.Join(parameters["username"], "")
	password := strings.Join(parameters["password"], "")

	return zeppelin.NewClient(hostname, username, password)
}

func list(w http.ResponseWriter, r *http.Request) {
	client, err := login(w, r)
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
		return
	}

	w.Write(b)
}

func run(w http.ResponseWriter, r *http.Request) {
	client, err := login(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not connect to %s: %v", hostname, err), http.StatusInternalServerError)
		return
	}

	parameters := r.URL.Query()
	notebookIDs := strings.Split(strings.Join(parameters["notebookIDs"], ""), ",")

	log.Println(notebookIDs)

	res, err := client.RunNotebooks(notebookIDs)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not run all notebooks: %v", err), http.StatusPartialContent)
	}

	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("could not format response from %s: %v", hostname, err), http.StatusInternalServerError)
		return
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
	srv.HandleFunc("/run", run)

	log.Printf("serving on %s", add)
	if err := http.ListenAndServe(add, srv); err != nil {
		log.Fatalf("server stopped: %v\n", err)
	}
}
