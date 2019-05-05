package main

import (
	"fmt"
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
	fmt.Printf("setting up target '%s'\n", *outDir)
	if err := os.MkdirAll(*outDir, os.ModePerm); err != nil {
		panic(err)
	}

	// generate index.html
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

	// generate non-template files
	for _, t := range []string{
		"app.js",
		"style.css",
	} {
		appData, err := internal.ReadFile("web/" + t)
		if err != nil {
			panic(err)
		}
		f, err = os.OpenFile(filepath.Join(*outDir, t), os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		if _, err := f.Write(appData); err != nil {
			panic(err)
		}
		f.Sync()
		f.Close()
	}

	fmt.Printf("generated web app in '%s'\n", *outDir)
}
