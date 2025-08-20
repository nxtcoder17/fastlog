package main

import (
	"html/template"
	"io"
)

var Template *template.Template = template.New("template:main")

type GetComponentFn func(attr map[string]any) (Component, error)

var Components map[string]GetComponentFn = make(map[string]GetComponentFn)

type Component interface {
	Render(w io.Writer) error
}
