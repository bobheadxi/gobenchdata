package web

import (
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

// ListenAndServe serves the web app on the given address
func ListenAndServe(addr string, config Config, it TemplateIndexHTML) error {
	if err := populateFileIndexHTML(it); err != nil {
		return err
	}

	// generate configuration in virtual filesystem
	f, err := FS.OpenFile(CTX, "/gobenchdata-web.yml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("failed to add config to virtual filesystem: %w", err)
	}
	b, _ := yaml.Marshal(&config)
	if _, err = f.Write(b); err != nil {
		return fmt.Errorf("failed to add config to virtual filesystem: %w", err)
	}
	if err = f.Close(); err != nil {
		return fmt.Errorf("failed to add config to virtual filesystem: %w", err)
	}

	// read benchmarks and add to virtual filesystem
	f, err = FS.OpenFile(CTX, "/benchmarks.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("failed to load benchmarks: %w", err)
	}
	b, err = os.ReadFile("./benchmarks.json")
	if err != nil {
		return fmt.Errorf("failed to load benchmarks: %w", err)
	}
	if _, err = f.Write(b); err != nil {
		return fmt.Errorf("failed to load benchmarks: %w", err)
	}
	if err = f.Close(); err != nil {
		return fmt.Errorf("failed to load benchmarks: %w", err)
	}

	// set up server
	handler := http.NewServeMux()
	handler.Handle("/", http.FileServer(HTTP))
	return http.ListenAndServe(addr, handler)
}
