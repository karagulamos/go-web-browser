package main

import (
	"fmt"

	"github.com/karagulamos/go-web-browser/internal/browser"
)

func main() {
	renderer := browser.NewRequest("https://example.org")

	result := renderer.Invoke()
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	fmt.Println(string(result.Content))
}
