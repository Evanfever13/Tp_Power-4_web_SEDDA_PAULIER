package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var listTemp *template.Template

func LoadTemplates() {
	listTemplate, errTemplate := template.ParseGlob("./../templates/*.html")
	if errTemplate != nil {
		log.Fatal("Erreur template:", errTemplate)
		return
	}
	listTemp = listTemplate
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {

	var buffer bytes.Buffer

	errRender := listTemp.ExecuteTemplate(&buffer, name, data)
	if errRender != nil {
		http.Redirect(w, r, fmt.Sprintf("/error?code=%d&message=%s",
			http.StatusInternalServerError,
			url.QueryEscape("Erreur lors du chargement de la page")),
			http.StatusSeeOther)
		return
	}
	buffer.WriteTo(w)
}
