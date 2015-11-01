package contollers

import (
	"fmt"
	"html/template"
	"io"
	"os"
	fpath "path/filepath"

	"github.com/revel/revel"
	"github.com/waiteb3/revel-swagger/modules/swaggify"
)

type Swaggify struct {
	*revel.Controller
}

type Template struct {
	*template.Template
	name string
}

func NewTemplate(name string) (*Template, error) {
	template, err := template.ParseFiles(fpath.Join(swaggify.ViewsPath, name))
	return &Template{template, name}, err
}

func (s *Template) Name() string {
	return s.name
}

func (s *Template) Content() []string {
	content, _ := revel.ReadLines(fpath.Join(swaggify.ViewsPath, s.name))
	return content
}

func (s *Template) Render(wr io.Writer, arg interface{}) error {
	return s.Execute(wr, arg)
}

var IndexTemplate *Template

func (c Swaggify) ServeUI(basePath string) revel.Result {
	// always recompiles in dev mode, once in nondev
	if revel.DevMode || IndexTemplate == nil {
		var err error
		IndexTemplate, err = NewTemplate("index.html")
		if err != nil {
			return c.RenderError(err)
		}
	}

	// TODO find better way to do this
	c.RenderArgs["SpecFile"] = c.Request.URL.Path + "swagger.json"
	return &revel.RenderTemplateResult{
		Template:   IndexTemplate,
		RenderArgs: c.RenderArgs,
	}
}

func (c Swaggify) ServeAssets(basePath, filepath string) revel.Result {
	if _, ok := swaggify.APIs[basePath]; !ok && false {
		return c.RenderError(fmt.Errorf("No matching API-UI at endpoint %s", basePath))
	}

	file, err := os.Open(fpath.Join(swaggify.AssetsPath, filepath))
	if err != nil {
		return c.RenderError(err)
	}

	// TODO look into inline vs attachment
	return c.RenderFile(file, revel.Inline)
}

func (c Swaggify) Spec(basePath string) revel.Result {
	if api := swaggify.APIs[basePath]; api != nil {
		return c.RenderJson(api)
	} else {
		return c.RenderError(fmt.Errorf("No matching API at endpoint %s", basePath))
	}
}
