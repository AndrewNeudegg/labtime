package templator

import (
	"fmt"

	"github.com/flosch/pongo2"

	log "github.com/sirupsen/logrus"
)

// CreateContext will generate a valid pongo2 context.
func CreateContext(contentMap map[string]interface{}) *pongo2.Context {
	context := pongo2.Context(contentMap)
	return &context
}

// RenderTemplate will apply a context to a jinja2 style template and return the result.
func RenderTemplate(templateFile string, context *pongo2.Context) (string, error) {
	log.Info(fmt.Sprintf("Loading template from file: %s", templateFile))
	tpl, err := pongo2.FromFile(templateFile)
	if err != nil {
		return "", err
	}

	out, err := tpl.Execute(*context)
	if err != nil {
		return "", err
	}
	log.Info(fmt.Sprintf("Rendered template"))

	return out, nil
}
