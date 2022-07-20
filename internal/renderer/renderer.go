package renderer

import "github.com/karagulamos/go-web-browser/pkg/network"

type Result struct {
	Headers network.Headers
	Content []byte
	Error   error
}

type Renderer interface {
	Invoke() *Result
}
