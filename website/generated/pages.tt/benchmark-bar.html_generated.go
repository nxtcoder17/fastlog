package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"html/template"
	"io"
	"strings"
)

func init() {
	Components["benchmarkbar"] = func(attr map[string]any) (Component, error) {
		return NewBenchmarkBar(attr)
	}

	parseBenchmarkBar()
}

func parseBenchmarkBar() error {
	_, err := Template.Parse(`{{- define "BenchmarkBar" }}
{{- /* @param name string */}}
{{- /* @param value string */}}
{{- /* @param color string (green|yellow|orange|red) */}}
{{- /* @param percentage int */}}
<div>
  <div class="flex justify-between mb-2">
    <span class="font-semibold">{{.name}}</span>
    <span class="text-{{.color}}-400 font-mono">{{.value}}</span>
  </div>
  <div class="bg-gray-700 rounded-full h-8 relative overflow-hidden">
    <div class="bg-{{.color}}-500 h-full"
         style="width: {{.percentage}}%">
    </div>
  </div>
</div>
{{- end }}
`)
	return err
}

type BenchmarkBar struct {
	Name       string `json:"name" validate:"required"`
	Color      string `json:"color" validate:"required"`
	Value      string `json:"value" validate:"required"`
	Percentage int    `json:"percentage" validate:"required"`

	// raw field contains all the
	// - known attributes (i.e. those defined above this line)
	// - and unkwnown ones (props) that are passed in html
	raw map[string]any `json:"-"`
}

func NewBenchmarkBar(attrs map[string]any) (*BenchmarkBar, error) {
	var s BenchmarkBar

	decoderCfg := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &s,
		TagName:          "json",
	}

	decoder, _ := mapstructure.NewDecoder(decoderCfg)
	if err := decoder.Decode(attrs); err != nil {
		panic(err)
	}

	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		fmt.Println("Validation failed:", err)
	}

	known := map[string]any{
		"name":       attrs["name"],
		"color":      attrs["color"],
		"value":      attrs["value"],
		"percentage": attrs["percentage"],
	}

	var unknown []string

	for k, v := range attrs {
		if _, ok := known[k]; !ok {
			unknown = append(unknown, fmt.Sprintf("%s=%q", k, v))
		}
	}

	s.raw = make(map[string]any, len(known)+1)
	for k, v := range known {
		s.raw[k] = v
	}

	_ = strings.ToLower

	s.raw["props"] = template.HTMLAttr(strings.Join(unknown, " "))

	return &s, nil
}

func (n *BenchmarkBar) TemplateName() string {
	return "BenchmarkBar"
}

func (n *BenchmarkBar) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(n)
}

func (n *BenchmarkBar) Render(w io.Writer) error {
	if err := n.Validate(); err != nil {
		return err
	}

	if n.raw == nil {
		b, err := json.Marshal(n)
		if err != nil {
			return err
		}

		n.raw = make(map[string]any)
		if err := json.Unmarshal(b, &n.raw); err != nil {
			return err
		}
	}
	return Template.ExecuteTemplate(w, n.TemplateName(), n.raw)
}
