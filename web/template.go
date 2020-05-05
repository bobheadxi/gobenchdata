package web

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

// TemplateIndexHTML declares variables for index.html
type TemplateIndexHTML struct {
	Title   string
	Headers []string
}

func populateFileIndexHTML(t TemplateIndexHTML) error {
	tmp, _ := template.New("index.html").Parse(string(FileIndexHTML))
	wr := bytes.NewBuffer([]byte{})
	if err := tmp.Execute(wr, struct {
		Title  string
		Header string
	}{
		Title:  t.Title,
		Header: strings.Join(t.Headers, "\n"),
	}); err != nil {
		return fmt.Errorf("failed to apply template to 'index.html': %w", err)
	}
	FileIndexHTML = wr.Bytes()
	return nil
}
