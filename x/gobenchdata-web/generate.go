package main

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/bobheadxi/gobenchdata/x/gobenchdata-web/internal"
)

type indexHTML struct {
	Title          string
	BenchmarksPath string
}

func generate() {
	// generate index.html
	os.MkdirAll(*outDir, os.ModePerm)
	tmpData, err := internal.ReadFile("web/index.html")
	if err != nil {
		panic(err)
	}
	tmp, err := template.New("index.html").Parse(string(tmpData))
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(filepath.Join(*outDir, "index.html"), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	if err := tmp.Execute(f, &indexHTML{
		Title:          *title,
		BenchmarksPath: *benchmarksPath,
	}); err != nil {
		panic(err)
	}
	f.Sync()
	f.Close()

	// generate app.js
	appData, err := internal.ReadFile("web/app.js")
	if err != nil {
		panic(err)
	}
	f, err = os.OpenFile(filepath.Join(*outDir, "app.js"), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	if _, err := f.Write(appData); err != nil {
		panic(err)
	}
	f.Sync()
	f.Close()
}
