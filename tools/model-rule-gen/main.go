package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"text/template"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/wata727/tflint/tools/utils"
)

var modelPathRoot = "rules/awsrules/models"
var mappingFile = fmt.Sprintf("%s/mapping.hcl", modelPathRoot)
var tmplFile = fmt.Sprintf("%s/pattern_rule.go.tmpl", modelPathRoot)

type mappings struct {
	Mapping []mapping `hcl:"mapping,block"`
}

type mapping struct {
	Resource resource `hcl:"resource,block"`
	Model    model    `hcl:"model,block"`
}

type resource struct {
	Type      string `hcl:"type"`
	Attribute string `hcl:"attribute"`
}

type model struct {
	Path  string `hcl:"path"`
	Shape string `hcl:"shape"`
}

type metadata struct {
	RuleName      string
	RuleNameCC    string
	ResourceType  string
	AttributeName string
	Max           int
	Min           int
	Pattern       string
}

func main() {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(mappingFile)
	if diags.HasErrors() {
		panic(diags)
	}

	var mappings mappings
	diags = gohcl.DecodeBody(f.Body, nil, &mappings)
	if diags.HasErrors() {
		panic(diags)
	}

	for _, mapping := range mappings.Mapping {
		raw, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", modelPathRoot, mapping.Model.Path))
		if err != nil {
			panic(err)
		}

		var api map[string]interface{}
		err = json.Unmarshal(raw, &api)
		if err != nil {
			panic(err)
		}
		shapes := api["shapes"].(map[string]interface{})

		resource := mapping.Resource.Type
		attribute := mapping.Resource.Attribute
		ruleName := fmt.Sprintf("%s_invalid_%s", resource, attribute)
		model := shapes[mapping.Model.Shape].(map[string]interface{})
		meta := &metadata{
			RuleName:      ruleName,
			RuleNameCC:    utils.ToCamel(ruleName),
			ResourceType:  resource,
			AttributeName: attribute,
			Max:           int(model["max"].(float64)),
			Min:           int(model["min"].(float64)),
			Pattern:       replacePattern(model["pattern"].(string)),
		}

		file, err := os.Create(fmt.Sprintf("%s/%s.go", modelPathRoot, ruleName))
		if err != nil {
			panic(err)
		}

		tmpl := template.Must(template.ParseFiles(tmplFile))
		err = tmpl.Execute(file, meta)
		if err != nil {
			panic(err)
		}
	}
}

func replacePattern(pattern string) string {
	reg := regexp.MustCompile(`\\u([0-9A-F]{4})`)
	return reg.ReplaceAllString(pattern, `\x{$1}`)
}
