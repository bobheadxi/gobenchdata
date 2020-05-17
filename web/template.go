package web

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
)

// TemplateIndexHTML declares variables for index.html
type TemplateIndexHTML struct {
	Title   string
	Headers []string
}

type internalIndexHTML struct {
	Title  string
	Header template.HTML
}

func populateFileIndexHTML(t TemplateIndexHTML) error {
	tmp, err := template.New("index.html").Parse(string(FileIndexHTML))
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}
	wr := bytes.NewBuffer([]byte{})
	if err := tmp.Execute(wr, &internalIndexHTML{
		Title:  t.Title,
		Header: template.HTML(strings.Join(t.Headers, "\n")),
	}); err != nil {
		return fmt.Errorf("failed to apply template to 'index.html': %w", err)
	}

	f, err := FS.OpenFile(CTX, "/index.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("failed to open template: %w", err)
	}
	if _, err = f.Write(wr.Bytes()); err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}
	if err = f.Close(); err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}
	return nil
}
