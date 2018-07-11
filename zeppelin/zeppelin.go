package zeppelin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
func (c *Client) ListNotebooks() (ListResponse, error) {
	res, err := c.login()
	if err != nil {
		return ListResponse{}, err
	}

	targetURL := urlWithPath("api/notebook", c.url)

	res, err = c.Get(targetURL.String())
	if err != nil {
		return ListResponse{}, fmt.Errorf("could not get %s: %v", targetURL.String(), err)
	}

	if res.StatusCode == 500 {
		return ListResponse{}, fmt.Errorf("remote service experiencing remote server error")
	}

	b, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var list ListResponse
	err = json.Unmarshal(b, &list)
	if err != nil {
		return ListResponse{}, fmt.Errorf("could not unmarshal response: %v", err)
	}

	return list, nil
}

func (c *Client) postrequest(url string, body io.Reader) (StdResponse, error) {
	_, err := c.login()
	if err != nil {
		return StdResponse{}, err
	}

	targetURL := urlWithPath(url, c.url)

	res, err := c.Post(targetURL.String(), "application/json", body)
	if err != nil {
		return StdResponse{}, fmt.Errorf("could not post to %s: %v", targetURL.String(), err)
	}

	if res.StatusCode == 500 {
		return StdResponse{}, fmt.Errorf("remote service experiencing remote server error: %v", res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return StdResponse{}, fmt.Errorf("could not read response: %v", err)
	}
	var response StdResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		return StdResponse{}, fmt.Errorf("could not unmarshal response: %v", err)
	}

	return response, nil
}

// NewNotebook creates a new notebook
func (c *Client) NewNotebook(newNote NewNoteRequestBody) (StdResponse, error) {
	b, err := json.Marshal(newNote)
	if err != nil {
		return StdResponse{}, fmt.Errorf("could not read input json: %v", err)
	}
	return c.postrequest("api/notebook", bytes.NewReader(b))
}

// RunNotebooks run notebooks in `notebookIDs []string`
func (c *Client) RunNotebooks(notebookIDs []string) ([]StdResponse, error) {

	responses := []StdResponse{}
	for _, notebookID := range notebookIDs {
		res, err := c.RunNotebook(notebookID)
		responses = append(responses, res)
		if err != nil {
			return responses, fmt.Errorf("could not run notebook %s: %v", notebookID, err)
		}
	}

	return responses, nil
}

// RunNotebook run notebook with ID `notebookID`
func (c *Client) RunNotebook(notebookID string) (StdResponse, error) {
	return c.postrequest(fmt.Sprintf("api/notebook/job/%s", notebookID), nil)
}

// GetNotePermission retrieves note permission for note `notebookID`
func (c *Client) GetNotePermission(notebookID string) (PermissionResponse, error) {
	res, err := c.login()
	if err != nil {
		return PermissionResponse{}, err
	}

	targetURL := urlWithPath(fmt.Sprintf("api/notebook/%s/permissions", notebookID), c.url)

	res, err = c.Get(targetURL.String())
	if err != nil {
		return PermissionResponse{}, fmt.Errorf("could not get %s: %v", targetURL.String(), err)
	}

	if res.StatusCode == 500 {
		return PermissionResponse{}, fmt.Errorf("remote service experiencing remote server error")
	}

	b, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var permissions PermissionResponse
	err = json.Unmarshal(b, &permissions)
	if err != nil {
		return PermissionResponse{}, fmt.Errorf("could not unmarshal response: %v", err)
	}

	return permissions, nil
}

// SetNotePermission set permission `permission` on notebook `notebookID`
func (c *Client) SetNotePermission(notebookID string, permission Permission) (StdResponse, error) {
	_, err := c.login()
	if err != nil {
		return StdResponse{}, err
	}

	targetURL := urlWithPath(fmt.Sprintf("api/notebook/%s/permissions", notebookID), c.url)

	b, err := json.Marshal(permission)
	if err != nil {
		return StdResponse{}, fmt.Errorf("could not marshal input json: %v", err)
	}

	request, err := http.NewRequest(http.MethodPut, targetURL.String(), bytes.NewReader(b))
	res, err := c.Do(request)
	if err != nil {
		return StdResponse{}, fmt.Errorf("could not put to %s: %v", targetURL.String(), err)
	}

	if res.StatusCode == 500 {
		return StdResponse{}, fmt.Errorf("remote service experiencing remote server error: %v", res.Status)
	}

	b, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return StdResponse{}, fmt.Errorf("could not read response: %v", err)
	}
	var response StdResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		return StdResponse{}, fmt.Errorf("could not unmarshal response: %v", err)
	}

	return response, nil
}
