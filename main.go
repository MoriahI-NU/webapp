package main

import (
	"fmt"
	"html/template"
	"net/http"

	// "webapp/assist"
	webapp "webapp/assist"
)

type Response struct {
	Title   string
	Content string
}

func ask(w http.ResponseWriter, r *http.Request) {
	var fileName = "ask.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	t.ExecuteTemplate(w, fileName, "What can I help you with?")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func askSubmit(w http.ResponseWriter, r *http.Request) {
	question := r.FormValue("question")
	fmt.Println(question)

	sectionsMap, roboticsHeaders := webapp.GatherInfo()
	appResponse := webapp.AppResponse(question, sectionsMap, roboticsHeaders)
	// fmt.Fprintf(w, "Topic: %s\n%s", question, appResponse)

	answer := Response{
		Title:   question,
		Content: appResponse,
	}

	var fileName2 = "submit.html"
	t2, err := template.ParseFiles(fileName2)
	if err != nil {
		fmt.Println(err)
		return
	}

	t2.ExecuteTemplate(w, fileName2, answer)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/ask":
		ask(w, r)
	case "/ask-submit":
		askSubmit(w, r)
	default:
		var fileName3 = "home.html"
		t3, err := template.ParseFiles(fileName3)
		if err != nil {
			fmt.Println(err)
			return
		}

		t3.ExecuteTemplate(w, fileName3, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe("", nil)
}
