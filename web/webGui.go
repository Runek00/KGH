package web

import (
	"html/template"
	"log"
	"net/http"
)

type Translations struct {
	FindInputLabel  string
	FindButtonLabel string
}

type PageData struct {
	Lang Translations
}

func WebGui() {
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/pullAll", pullAllHandler)
	http.HandleFunc("/pullEvents", pullEvents)

	log.Fatal(http.ListenAndServe(":8079", nil))
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/templates/index.html"))
	tmpl.Execute(w, getMainPageData())
}

func getMainPageData() PageData {
	return PageData{Translations{"input", "find"}}
}
