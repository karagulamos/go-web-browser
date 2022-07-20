package renderer

import "strings"

type DataRenderer struct {
	url string
}

func NewDataRenderer(url string) Renderer {
	url = strings.TrimPrefix(url, "data:")

	return NewHtmlRenderer(&DataRenderer{url}) // ASSUMPTION: data:text/html
}

func (r *DataRenderer) Invoke() *Result {
	parts := strings.SplitN(r.url, ",", 2)
	if len(parts) != 2 {
		return &Result{Error: ErrInvalidDataURI}
	}

	contentType, data := parts[0], parts[1]

	if !strings.HasPrefix(contentType, "text/html") { // ASSUMPTION: data:text/html
		return &Result{Error: ErrUnsupportedContentType}
	}

	body := "<html><body>" + data + "</body></html>"
	return &Result{Content: []byte(body)}
}
