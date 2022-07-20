package browser

import (
	"strings"

	"github.com/karagulamos/go-web-browser/internal/renderer"
)

const MaxHttpRedirects = 5

func NewRequest(url string) renderer.Renderer {
	switch url = strings.TrimSpace(url); {
	case strings.HasPrefix(url, "http://"):
		fallthrough
	case strings.HasPrefix(url, "https://"):
		return renderer.NewHtmlRenderer(
			renderer.NewSiteRenderer(url, MaxHttpRedirects))
	case strings.HasPrefix(url, "file:///"):
		return renderer.NewFSRenderer(url)
	case strings.HasPrefix(url, "view-source:"):
		return renderer.NewSiteRenderer(
			strings.TrimPrefix(url, "view-source:"), MaxHttpRedirects)
	case strings.HasPrefix(url, "data:"):
		return renderer.NewDataRenderer(url)
	default:
		return renderer.NewNoopRenderer()
	}
}
