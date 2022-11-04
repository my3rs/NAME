package model

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"html/template"
)

func Markdown2Html(markdown string) template.HTML {
	unsafe := blackfriday.Run([]byte(markdown))
	html := template.HTML(bluemonday.UGCPolicy().SanitizeBytes(unsafe))

	return html
}
