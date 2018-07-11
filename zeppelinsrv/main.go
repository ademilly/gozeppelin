package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ademilly/gozeppelin/zeppelin"
)

var client zeppelin.Client
var hostname string
var port string
var blocked bool
var rate time.Duration

func usage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf(`
endpoints:
  - list => list notebooks available on %s
  - run  => run notebooks by IDs given as comma separated URL parameter notebookIDs
    `, hostname)))
}

func newClient(_ http.ResponseWriter, r *http.Request) (*zeppelin.Client, error) {
	parameters := r.URL.Query()

	username := strings.Join(parameters["username"], "")
	password := strings.Join(parameters["password"], "")

	return zeppelin.NewClient(hostname, username, password)
}

func list(w http.ResponseWriter, r *http.Request) {
	client, err := newClient(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not connect to %s: %v", hostname, err), http.StatusInternalServerError)
		return
	}

	log.Printf("listing notebooks request from %s", r.RemoteAddr)
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
	client, err := newClient(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not connect to %s: %v", hostname, err), http.StatusInternalServerError)
		return
	}

	if blocked {
		w.Write([]byte("computation still going on, try again later ;)"))
		return
	}

	blocked = true
	go func() {
		time.AfterFunc(rate, func() {
			blocked = false
		})
	}()

	parameters := r.URL.Query()
	notebookIDs := strings.Split(strings.Join(parameters["notebookIDs"], ""), ",")

	log.Printf("runnings notebooks %v request from %s", notebookIDs, r.RemoteAddr)
	res, err := client.RunNotebooks(notebookIDs)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not run all notebooks: %v", err), http.StatusPartialContent)
		return
	}

	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("could not format response from %s: %v", hostname, err), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func init() {
	blocked = false
	rate = 10 * time.Minute
}

func main() {
	flag.StringVar(&port, "port", "8080", "port number on which to serve")
	flag.StringVar(&hostname, "hostname", "localhost", "zeppelin service hostname")
	flag.Parse()

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	handler := http.NewServeMux()

	handler.HandleFunc("/", usage)
	handler.HandleFunc("/list", list)
	handler.HandleFunc("/run", run)

	srv := &http.Server{
		Addr:        addr,
		ReadTimeout: 5 * time.Second,
		Handler:     http.TimeoutHandler(handler, 1*time.Minute, "server not available for now, try later ;)"),
	}

	log.Printf("serving on %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server stopped: %v\n", err)
	}
}
