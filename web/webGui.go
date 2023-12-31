package web

import (
	"KGH/base"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/atotto/clipboard"
)

type Translations struct {
	FindInputLabel  string
	FindButtonLabel string
	Settings        string
	DefaultTemplate string
	FilePath        string
}

type PageData struct {
	Lang Translations
}

func WebGui() {
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/pullAll", pullAllHandler)
	http.HandleFunc("/pullEvents", pullEvents)
	http.HandleFunc("/find", find)
	http.HandleFunc("/settings", settings)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		tmpl := template.Must(template.ParseFiles("static/templates/index.html"))
		tmpl.Execute(w, getMainPageData())
	} else {
		handler := http.FileServer(http.Dir("./static"))
		handler.ServeHTTP(w, r)

	}
}

func getMainPageData() PageData {
	return PageData{Translations{"input", "find", "Settings", "Default template", "Output file path"}}
}

func find(w http.ResponseWriter, r *http.Request) {
	found := base.FindCommits(r.FormValue("input"))
	output := ""
	for _, commit := range found {
		output += "<p>" + commit + "</p>\n"
	}
	copiable := strings.Join(found, "\n")
	clipboard.WriteAll(copiable)
	w.Write([]byte(output))
}

func settings(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("static/templates/settings.html"))
		tmpl.Execute(w, getMainPageData())
	}
}
