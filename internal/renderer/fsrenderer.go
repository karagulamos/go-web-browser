package renderer

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

type FSRenderer struct {
	path string
}

func NewFSRenderer(url string) Renderer {
	return &FSRenderer{
		path: strings.TrimPrefix(url, "file://"),
	}
}

func (r *FSRenderer) Invoke() *Result {
	fi, err := os.Stat(r.path)
	if os.IsNotExist(err) {
		return &Result{Error: err}
	}

	if !fi.IsDir() {
		content, err := ioutil.ReadFile(r.path)
		if err != nil {
			return &Result{Error: err}
		}

		return &Result{Content: content}
	}

	var output bytes.Buffer

	if fi, err := ioutil.ReadDir(r.path); err == nil {
		for _, f := range fi {
			fileName := f.Name()

			if f.IsDir() {
				fileName += "/"
			}

			output.WriteString(fileName)
			output.WriteString("\n")
		}
	}

	return &Result{Content: output.Bytes()}
}
