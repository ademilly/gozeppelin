package zeppelin

// TextToNewNote generates a NewNoteRequestBody from name, title and text strings
func TextToNewNote(name, title, text string) NewNoteRequestBody {
	return NewNoteRequestBody{
		name,
		[]struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{
			{title, text},
		},
	}
}
