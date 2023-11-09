package main

import (
	"html/template"
	"log"
	"net/http"
)

func webGui() {
	http.HandleFunc("/", mainPage)

	log.Fatal(http.ListenAndServe(":8079", nil))
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/templates/index.html"))
	tmpl.Execute(w, nil)
}
