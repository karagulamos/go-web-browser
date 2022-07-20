package renderer

import "fmt"

var (
	ErrNotImplemented             = fmt.Errorf("not implemented")
	ErrTooManyRedirects           = fmt.Errorf("too many redirects")
	ErrInvalidSiteURL             = fmt.Errorf("invalid site URL")
	ErrInvalidDataURI             = fmt.Errorf("invalid data URI")
	ErrUnsupportedContentType     = fmt.Errorf("unsupported content type")
	ErrUnsupportedContentEncoding = fmt.Errorf("unsupported content encoding")
)
