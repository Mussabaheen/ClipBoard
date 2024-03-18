package main

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"text/template"

	"golang.design/x/clipboard"
)

func main() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)

	go copyFromClipBoard(ch)
	http.HandleFunc("/", ShowBooks)
	http.ListenAndServe(":8080", nil)
}

func copyFromClipBoard(ch <-chan []byte) {
	for data_from_clipboard := range ch {
		data = append(data, string(data_from_clipboard))
		fmt.Println(string(data_from_clipboard))
		fmt.Println("***************************************************************************************************************************")
	}
}

func ShowBooks(w http.ResponseWriter, r *http.Request) {
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

var data []string
