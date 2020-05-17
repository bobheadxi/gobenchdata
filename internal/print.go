package internal

import "fmt"

// Printer only outputs if debug=true
type Printer struct {
	Debug bool
}

// Print wraps fmt.Print
func (p *Printer) Print(v ...interface{}) (int, error) {
	if p.Debug {
		return fmt.Print(v...)
	}
	return 0, nil
}

// Printf wraps fmt.Printf
func (p *Printer) Printf(f string, v ...interface{}) (int, error) {
	if p.Debug {
		return fmt.Printf(f+"\n", v...)
	}
	return 0, nil
}
