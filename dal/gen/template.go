package gen

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"github.com/winjo/codegen/dal/util"
)

//go:embed tmpl/*
var tmplfs embed.FS

type Template struct {
	*template.Template
}

var tmpl *Template

func init() {
	t, err := template.New("").ParseFS(tmplfs, "tmpl/*")
	util.AssertNotNil(err)
	tmpl = &Template{
		Template: t,
	}
}

func (t *Template) Exec(name string, data map[string]any) string {
	buf := new(bytes.Buffer)
	err := t.ExecuteTemplate(buf, fmt.Sprintf("%s.tmpl", name), data)
	util.AssertNotNil(err)
	return buf.String()
}
