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

// Permission struct holds notebook permission data
type Permission struct {
	Owners  []string `json:"owners"`
	Readers []string `json:"readers"`
	Writers []string `json:"writers"`
}

// PermissionResponse struct holds response to get a note permission on zeppelin
type PermissionResponse struct {
	Status  string
	Message string
	Body    Permission
}

// StdResponse struct holds standard zeppelin response
type StdResponse struct {
	Status  string
	Message string `json:",omitempty"`
	Body    string `json:",omitempty"`
}
