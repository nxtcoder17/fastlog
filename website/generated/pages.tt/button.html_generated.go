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
	Components["styledbutton"] = func(attr map[string]any) (Component, error) {
		return NewStyledButton(attr)
	}

	parseButton()
}

func parseButton() error {
	_, err := Template.Parse(`{{- define "StyledButton" }}
{{- /* @param href string */}}
{{- /* @param variant string (primary|secondary) */}}
{{- $classes := "" }}
{{- if eq .variant "primary" }}
{{- $classes = "bg-purple-600 hover:bg-purple-700 px-8 py-3 rounded-lg font-semibold transition" }}
{{- else if eq .variant "secondary" }}
{{- $classes = "border border-purple-600 hover:bg-purple-600/10 px-8 py-3 rounded-lg font-semibold transition" }}
{{- end }}
<a href="{{.href}}"
   class="{{$classes}}">
  <Children />
</a>
{{- end }}
`)
	return err
}

type StyledButton struct {
	Variant string `json:"variant" validate:"required"`
	Href    string `json:"href" validate:"required"`

	// raw field contains all the
	// - known attributes (i.e. those defined above this line)
	// - and unkwnown ones (props) that are passed in html
	raw map[string]any `json:"-"`
}

func NewStyledButton(attrs map[string]any) (*StyledButton, error) {
	var s StyledButton

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
		"variant": attrs["variant"],
		"href":    attrs["href"],
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

func (n *StyledButton) TemplateName() string {
	return "StyledButton"
}

func (n *StyledButton) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(n)
}

func (n *StyledButton) Render(w io.Writer) error {
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
