package zeppelin

// NewNoteRequestBody struct represents a new note request body as in
// https://zeppelin.apache.org/docs/latest/rest-api/rest-notebook.html#create-a-new-note
type NewNoteRequestBody struct {
	Name       string `json:"name"`
	Paragraphs []struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"paragraphs"`
}
