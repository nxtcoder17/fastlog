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
	Components["codeblock"] = func(attr map[string]any) (Component, error) {
		return NewCodeBlock(attr)
	}

	parseCodeBlock()
}

func parseCodeBlock() error {
	_, err := Template.Parse(`{{- define "CodeBlock" }}
{{- /* @param language? string */}}
<pre class="bg-gray-800 rounded-lg p-6 text-left overflow-x-auto">
  <code class="text-sm">
  <Children />
  </code>
</pre>
{{- end }}
`)
	return err
}

type CodeBlock struct {

	// raw field contains all the
	// - known attributes (i.e. those defined above this line)
	// - and unkwnown ones (props) that are passed in html
	raw map[string]any `json:"-"`
}

func NewCodeBlock(attrs map[string]any) (*CodeBlock, error) {
	var s CodeBlock

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

	known := map[string]any{}

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

func (n *CodeBlock) TemplateName() string {
	return "CodeBlock"
}

func (n *CodeBlock) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(n)
}

func (n *CodeBlock) Render(w io.Writer) error {
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
