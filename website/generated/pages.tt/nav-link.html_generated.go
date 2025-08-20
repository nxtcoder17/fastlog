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
	Components["navlink"] = func(attr map[string]any) (Component, error) {
		return NewNavLink(attr)
	}

	parseNavLink()
}

func parseNavLink() error {
	_, err := Template.Parse(`{{- define "NavLink" }}
{{- /* @param href string */}}
<a href="{{.href}}"
   class="hover:text-purple-400 transition">
  <Children />
</a>
{{- end }}
`)
	return err
}

type NavLink struct {
	Href string `json:"href" validate:"required"`

	// raw field contains all the
	// - known attributes (i.e. those defined above this line)
	// - and unkwnown ones (props) that are passed in html
	raw map[string]any `json:"-"`
}

func NewNavLink(attrs map[string]any) (*NavLink, error) {
	var s NavLink

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
		"href": attrs["href"],
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

func (n *NavLink) TemplateName() string {
	return "NavLink"
}

func (n *NavLink) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(n)
}

func (n *NavLink) Render(w io.Writer) error {
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
