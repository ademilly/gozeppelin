package zeppelin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

// NewNoteRequestBody struct represents a new note request body as in
// https://zeppelin.apache.org/docs/latest/rest-api/rest-notebook.html#create-a-new-note
type NewNoteRequestBody struct {
	Name       string `json:"name"`
	Paragraphs []struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"paragraphs"`
}

// Client struct represents http client for Zeppelin
type Client struct {
	*http.Client
	user struct {
		username string
		password string
	}
	url *url.URL
}

func urlWithPath(path string, url *url.URL) url.URL {
	newURL := *url
	newURL.Path = path
	return newURL
}

// NewClient builds a new Zeppelin client
func NewClient(hostname, username, password string) (*Client, error) {
	URL, err := url.Parse(hostname)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	return &Client{
		&http.Client{Jar: jar},
		struct {
			username string
			password string
		}{username, password},
		URL,
	}, nil
}

func (c *Client) login() (*http.Response, error) {
	targetURL := urlWithPath("api/login", c.url)

	res, err := c.PostForm(targetURL.String(), url.Values{"userName": {c.user.username}, "password": {c.user.password}})
	if err != nil {
		return nil, fmt.Errorf("could not login to %s: %v", targetURL.String(), err)
	}

	return res, nil
}

// ListNotebooks lists notebooks
func (c *Client) ListNotebooks() (*http.Response, error) {
	res, err := c.login()
	if err != nil {
		return nil, err
	}

	targetURL := urlWithPath("api/notebook", c.url)

	res, err = c.Get(targetURL.String())
	if err != nil {
		return nil, fmt.Errorf("could not get %s: %v", targetURL.String(), err)
	}

	if res.StatusCode == 500 {
		return nil, fmt.Errorf("remote service experiencing remote server error")
	}

	return res, nil
}

// NewNotebook creates a new notebook
func (c *Client) NewNotebook(newNote NewNoteRequestBody) (*http.Response, error) {
	res, err := c.login()
	if err != nil {
		return nil, err
	}

	targetURL := urlWithPath("api/notebook", c.url)

	b, err := json.Marshal(newNote)
	res, err = c.Post(targetURL.String(), "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("could not post to %s: %v", targetURL.String(), err)
	}

	if res.StatusCode == 500 {
		return nil, fmt.Errorf("remote service experiencing remote server error")
	}

	return res, nil
}
