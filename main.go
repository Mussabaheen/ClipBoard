package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"path"
	"text/template"

	"golang.design/x/clipboard"
)

var data []string
var clients = make(map[chan []byte]struct{})

var (
	localPortUsage = "Specify the port number on which you want to serve the UI, by default 8080"
)

func main() {
	// localPort represents the arg with flag -p
	localPort := flag.String("port", "8080", localPortUsage)
	flag.Parse()
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)

	go copyFromClipBoard(ch)

	http.HandleFunc("/", ShowClipboard)
	http.HandleFunc("/updates", UpdatesHandler)
	err = http.ListenAndServe(":"+*localPort, nil)
	if err != nil {
		panic("error occured while running clipboard : " + err.Error())
	}
}

func copyFromClipBoard(ch <-chan []byte) {
	for data_from_clipboard := range ch {
		newData := string(data_from_clipboard)
		data = append([]string{newData}, data...)

		// Serialize data to JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshaling data to JSON:", err)
			continue
		}

		// Send update to clients
		for client := range clients {
			client <- jsonData
		}
	}
}

func ShowClipboard(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("internal/templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UpdatesHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a channel for this client
	client := make(chan []byte)
	// Register this client
	clients[client] = struct{}{}
	defer func() {
		// Unregister this client
		delete(clients, client)
		close(client)
	}()

	// Listen for updates from the channel
	for jsonData := range client {
		// Send update to the client
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		w.(http.Flusher).Flush()
	}
}
