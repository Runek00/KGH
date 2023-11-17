package main

import (
	"fmt"
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
	tmpl.Parse(`<div id="yolo" name="sse" hx-ext="sse" sse-connect="/pullEvents" sse-swap="message" hx-swap="afterend">`)
	tmpl.Execute(w, nil)
}

func pullEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")
	statusChan := make(chan string)
	defer func() {
		close(statusChan)
		statusChan = nil
	}()
	go func(flusher http.Flusher) {
		for {
			if statusChan == nil {
				break
			}
			fmt.Fprintf(w, "type: message\ndata: "+<-statusChan)
			flusher.Flush()
		}
		fmt.Fprintf(w, "type: finished")
		flusher.Flush()
	}(w.(http.Flusher))
	PullAll(statusChan)
}
