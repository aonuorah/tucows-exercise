package service

import (
	"bytes"
	"html/template"

	"github.com/aonuorah/tucows-exercise/internal/photo"
	"github.com/aonuorah/tucows-exercise/internal/quote"
	"github.com/spf13/afero"
)

type tmpl struct {
	filename string
	template *template.Template
}

func (t *tmpl) Template() *template.Template {
	if t.template == nil {
		content, _ := afero.ReadFile(FS, t.filename)
		t.template = template.Must(template.New("html").Parse(string(content)))
	}
	return t.template
}

var htmltmpl = &tmpl{filename: "web/tmpl/index.html.tmpl"}

// Builds the html for the given photo and quote
func BuildHTML(p *photo.Info, q *quote.Info) ([]byte, error) {
	tmpldata := map[string]any{
		"Category":    q.Category,
		"QuoteText":   q.Text,
		"Grayscale":   p.Grayscale,
		"ImageBase64": p.Base64Encoding(),
		"ImageWidth":  p.Width,
		"ImageHeight": p.Height,
	}

	htmlbuf := bytes.Buffer{}
	if err := htmltmpl.Template().Execute(&htmlbuf, tmpldata); err != nil {
		return nil, err
	}

	return htmlbuf.Bytes(), nil
}
