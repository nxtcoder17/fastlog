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
	Components["featurecard"] = func(attr map[string]any) (Component, error) {
		return NewFeatureCard(attr)
	}

	parseFeatureCard()
}

func parseFeatureCard() error {
	_, err := Template.Parse(`{{- define "FeatureCard" }}
{{- /* @param icon string */}}
{{- /* @param title string */}}
<div class="feature-card p-8 rounded-xl">
  <div class="text-4xl mb-4">
    {{.icon}}
  </div>
  <h4 class="text-2xl font-bold mb-3">
    {{.title}}
  </h4>
  <p class="text-gray-400">
    <Children />
  </p>
</div>
{{- end }}

`)
	return err
}

type FeatureCard struct {
	Icon  string `json:"icon" validate:"required"`
	Title string `json:"title" validate:"required"`

	// raw field contains all the
	// - known attributes (i.e. those defined above this line)
	// - and unkwnown ones (props) that are passed in html
	raw map[string]any `json:"-"`
}

func NewFeatureCard(attrs map[string]any) (*FeatureCard, error) {
	var s FeatureCard

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
		"icon":  attrs["icon"],
		"title": attrs["title"],
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

func (n *FeatureCard) TemplateName() string {
	return "FeatureCard"
}

func (n *FeatureCard) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(n)
}

func (n *FeatureCard) Render(w io.Writer) error {
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
