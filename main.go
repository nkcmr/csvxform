package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"text/template"
)

type tmplData struct {
	headers map[string]int
	row     []string
}

func (t tmplData) Col(s string) (string, error) {
	idx, ok := t.headers[s]
	if !ok {
		return "", fmt.Errorf("unknown column: '%s'", s)
	}
	return t.row[idx], nil
}

func _main() error {
	r := csv.NewReader(os.Stdin)
	header := true
	dat := tmplData{}
	if len(os.Args) == 1 {
		return errors.New("no template provided")
	}
	tmpl, err := template.New("row_tmpl").Parse(os.Args[1])
	if err != nil {
		return fmt.Errorf("failed to parse template: %s", err.Error())
	}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("failed to read stdin: %s", err.Error())
		}
		if header {
			header = false
			dat.headers = make(map[string]int, len(row))
			for idx, h := range row {
				dat.headers[h] = idx
			}
			continue
		}
		dat.row = row
		err = tmpl.Execute(os.Stdout, dat)
		if err != nil {
			return fmt.Errorf("failed to execute template: %s", err.Error())
		}
		fmt.Fprintln(os.Stdout)
	}
	return nil
}

func main() {
	if err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err.Error())
		os.Exit(1)
	}
}
