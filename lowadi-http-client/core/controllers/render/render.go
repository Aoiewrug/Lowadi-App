package render

import (
	"fmt"
	"html/template"
	"net/http"
)

var tempFolder = "C:/Users/3/Documents/GitHub/http-rest-client/front/"

// One way to display 1 page
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	parsedTemplate, err := template.ParseFiles(tempFolder + tmpl)
	if err != nil {
		fmt.Println("can't find the template", err)
	}

	x := parsedTemplate.ExecuteTemplate(w, tmpl, data)

	if x != nil {
		fmt.Println("error parsing the template", x)
	}
}

var templates = template.Must(template.ParseGlob("C:/Users/3/Documents/GitHub/http-rest-client/front/*"))

// Second way to display multiple pages
func RenderMultipleTemplates(w http.ResponseWriter, defineName string, data interface{}) {
	// you access the cached templates with the defined name, not the filename
	err := templates.ExecuteTemplate(w, defineName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
