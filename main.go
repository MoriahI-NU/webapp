package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func ask(w http.ResponseWriter, r *http.Request) {
	var fileName = "ask.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.ExecuteTemplate(w, fileName, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func askSubmit(w http.ResponseWriter, r *http.Request) {
	question := r.FormValue("question")

}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/ask":
		ask(w, r)
	case "/ask-submit":
	default:
		fmt.Fprintf(w, "Hello")
	}
}
func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("", nil)
}
