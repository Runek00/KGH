package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/r3labs/sse"
)

type Translations struct {
	FindInputLabel  string
	FindButtonLabel string
}

type PageData struct {
	Lang Translations
}

func webGui() {
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/pullAll", pullAllHandler)
	http.HandleFunc("/pullEvents", pullEvents)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/templates/index.html"))
	tmpl.Execute(w, getMainPageData())
}

func getMainPageData() PageData {
	return PageData{Translations{"input", "find"}}
}

func pullAllHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("sse-element")
	tmpl.Parse(`<div id="yolo" name="sse" hx-ext="sse" sse-connect="/pullEvents?stream=message" sse-swap="message, finished" hx-swap="afterend" _="on finished remove me">`)
	tmpl.Execute(w, nil)
}

func pullEvents(w http.ResponseWriter, r *http.Request) {
	server := sse.New()
	server.CreateStream("message")
	go server.HTTPHandler(w, r)
	statusChan := make(chan string)
	defer func() {
		close(statusChan)
		statusChan = nil
	}()
	go func() {
		for status := range statusChan {
			server.Publish("message", &sse.Event{
				Data:  []byte("<p>" + status + "</p>"),
				Event: []byte("message"),
			})
		}
		server.Publish("message", &sse.Event{
			Event: []byte("finished"),
		})
		server.Close()
	}()
	PullAll(statusChan)
}
