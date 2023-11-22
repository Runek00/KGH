package web

import (
	"KGH/base"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/r3labs/sse"
)

func pullAllHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("sse-element")
	tmpl.Parse(`<div id="yolo" name="sse" hx-ext="sse" sse-connect="/pullEvents?stream=message" sse-swap="message" hx-swap="afterend" _="on load add @disabled to #pullButton">`)
	tmpl.Execute(w, nil)
}

func pullEvents(w http.ResponseWriter, r *http.Request) {
	server := sse.New()
	server.CreateStream("message")
	go server.HTTPHandler(w, r)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		statusChan := make(chan string)

		wg2 := sync.WaitGroup{}
		wg2.Add(1)
		go func() {
			for status := range statusChan {
				server.Publish("message", &sse.Event{
					Data:  []byte("<p>" + status + "</p>"),
					Event: []byte("message"),
				})
			}
			wg2.Done()
		}()
		base.PullAll(statusChan)
		go func() {
			defer close(statusChan)
			statusChan <- "<div _='on load remove @disabled from pullButton then remove #yolo then remove me'></div>"
			time.Sleep(time.Second)
		}()
		wg2.Wait()
		statusChan = nil
		server.Close()
		wg.Done()
	}()
	wg.Wait()
}
