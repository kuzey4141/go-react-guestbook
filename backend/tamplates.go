package main

import (
	"html/template"
	"log"
)

// NewTemplateCache parses the given HTML file and returns a template.
func NewTemplateCache(file string) *template.Template {
	tmpl, err := template.ParseFiles(file)
	if err != nil {
		log.Fatalf("template %q could not be loaded: %v", file, err)
	}
	return tmpl
}
