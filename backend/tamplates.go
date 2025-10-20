package main

import (
	"html/template"
	"log"
)

// NewTemplateCache, belirtilen HTML dosyasını ayrıştırır ve bir şablon nesnesi döndürür.
func NewTemplateCache(file string) *template.Template {
	// template.Must, şablon okunurken hata olursa programı durdurur (Fatal).
	tmpl, err := template.ParseFiles(file)
	if err != nil {
		// log.Fatal, programı anında durdurur. Şablon olmadan çalışamayız.
		log.Fatalf("Şablon '%s' yüklenemedi: %v", file, err)
	}
	return tmpl
}
