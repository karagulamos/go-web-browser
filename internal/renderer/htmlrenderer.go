package renderer

import (
	"bytes"
	"unicode"
)

var (
	htmlEntities = map[string]string{
		"&amp;":   "&",
		"&lt;":    "<",
		"&gt;":    ">",
		"&quot;":  `"`,
		"&apos;":  `'`,
		"&nbsp;":  " ",
		"&copy;":  "(c)",
		"&ndash;": "-",
	}
)

type HtmlRenderer struct {
	source Renderer
}

func NewHtmlRenderer(source Renderer) Renderer {
	return &HtmlRenderer{
		source: source,
	}
}

func (r *HtmlRenderer) Invoke() *Result {
	result := r.source.Invoke()
	if result.Error != nil {
		return result
	}

	var output, htmlEntity bytes.Buffer
	var inAngle, inBody bool
	var element string

	for _, c := range result.Content {
		if c == '<' {
			inAngle = true
		} else if c == '>' {
			inAngle = false

			if element == "body" {
				inBody = true
			} else if element == "/body" {
				inBody = false
			}
			element = ""
		} else if !inAngle && inBody {
			// Replace HTML entities with their symbolic representation.
			// Did we find a HTML entity? Or are we in the middle of an HTML entity?
			if c == '&' || htmlEntity.Len() > 0 {
				// Append the current character to the HTML entity
				htmlEntity.WriteByte(c)

				// Check if we're at the end of the HTML entity
				if c != ';' {
					continue
				}

				entity := htmlEntity.String()

				// Check if the entity is valid and replace it with its symbolic representation
				// Otherwise, just write the entity to the output buffer
				if symbol := htmlEntities[entity]; symbol != "" {
					output.WriteString(symbol)
				} else {
					output.WriteString(entity)
				}

				htmlEntity.Reset()
			} else {
				output.WriteByte(c) // Write the current character to the output buffer
			}
		} else if inAngle {
			if !unicode.IsSpace(rune(c)) {
				element += string(c)
			}
		}
	}

	return &Result{Headers: result.Headers, Content: output.Bytes()}
}
