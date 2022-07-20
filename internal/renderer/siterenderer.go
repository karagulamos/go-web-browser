package renderer

import (
	"io/ioutil"
	"net"
	"strings"

	"github.com/karagulamos/go-web-browser/pkg/network"
)

type SiteRenderer struct {
	url       string
	redirects int
	socket    network.Socket
}

func NewSiteRenderer(url string, redirects int) Renderer {
	return &SiteRenderer{
		url, redirects, network.NewSocket(
			network.Tcp,
			network.Config{Secure: strings.HasPrefix(url, "https://")},
		),
	}
}

func (r *SiteRenderer) Invoke() *Result {
	if r.redirects == 0 {
		return &Result{Error: ErrTooManyRedirects}
	}

	parts := strings.SplitN(r.url, "://", 2)
	if len(parts) == 0 {
		return &Result{Error: ErrInvalidSiteURL}
	}

	scheme, host, path, port := parts[0], parts[1], "/", "80"

	if scheme == "https" {
		port = "443"
	}

	parts = strings.SplitN(host, "/", 2)
	if len(parts) == 2 {
		host, path = parts[0], "/"+parts[1]
	}

	parts = strings.SplitN(host, ":", 2)
	if len(parts) == 2 {
		host, port = parts[0], parts[1]
	}

	if err := r.socket.Connect(net.JoinHostPort(host, port)); err != nil {
		return &Result{Error: err}
	}

	r.socket.WritefLine("GET %s HTTP/1.1", path)
	r.socket.WritefLine("Host: %s", host)
	r.socket.WritefLine("Connection: close")
	r.socket.WritefLine("User-Agent: %s", "Mozilla/5.0")
	r.socket.WritefLine("")

	response, err := network.AsHttpResponse(r.socket)
	if err != nil {
		return &Result{Error: err}
	}
	defer response.Body.Close()

	if location := response.Headers.Get("location"); location != "" {
		return NewSiteRenderer(location, r.redirects-1).Invoke()
	}

	html, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &Result{Error: err}
	}

	return &Result{Headers: response.Headers, Content: html}
}
