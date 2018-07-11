package zeppelin

// notebook struct represents minimal information on Notebook object
type notebook struct {
	Name string
	ID   string
}

// ListResponse struct holds response to list query on zeppelin
type ListResponse struct {
	Status  string
	Message string
	Body    []notebook
}

// StdResponse struct holds standard zeppelin response
type StdResponse struct {
	Status  string
	Message string `json:",omitempty"`
	Body    string `json:",omitempty"`
}
