package zeppelin

// Notebook struct represents minimal information on Notebook object
type Notebook struct {
	Name string
	ID   string
}

// ListResponse struct holds response to list query on zeppelin
type ListResponse struct {
	Status  string
	Message string
	Body    []Notebook
}
