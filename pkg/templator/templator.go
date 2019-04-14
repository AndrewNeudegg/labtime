package templator

import (
	"io/ioutil"

	"github.com/flosch/pongo2"
	"gopkg.in/yaml.v2"
)

// CreateContext will generate a valid pongo2 context.
func CreateContext(contentMap map[string]interface{}) *pongo2.Context {
	context := pongo2.Context(contentMap)
	return &context
}

// RenderTemplate will apply a context to a jinja2 style template and return the result.
func RenderTemplate(templateFile string, context *pongo2.Context) (string, error) {
	tpl, err := pongo2.FromFile(templateFile)
	if err != nil {
		return "", err
	}

	out, err := tpl.Execute(*context)
	if err != nil {
		return "", err
	}

	return out, nil
}

// ReadYaml is a helper method for reading arbitrary yaml files into a pongo2.Context.
func ReadYaml(yamlFile string) (pongo2.Context, error) {
	var ret pongo2.Context

	yamlFileBytes, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFileBytes, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
