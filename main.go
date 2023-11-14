package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type Quote struct {
	Content    string `json:"content"`
	Author     string `json:"author"`
	AuthorSlug string `json:"authorSlug"`
}

func getQuote() (Quote, error) {
	resp, err := http.Get("https://api.quotable.io/random")
	if err != nil {
		return Quote{}, err
	}
	defer resp.Body.Close()

	var quote Quote

	err = json.NewDecoder(resp.Body).Decode(&quote)
	if err != nil {
		return Quote{}, err
	}
	return quote, nil

}

func handleQuote(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	quote, err := getQuote()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, quote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func generateQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	quote, err := getQuote()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	generate := fmt.Sprintf("<p id='quote'> <q> %s </q> - %s </p>", quote.Content, quote.Author)

	w.Header().Set("Content-Type", "application/json")
	//err = json.NewEncoder(w).Encode(map[string]string{"content": quote.Content, "author": quote.Author, "authorSlug": quote.AuthorSlug})

	t, _ := template.New("t").Parse(generate)

	// t.ExecuteTemplate(w, "film-list-element", Book{Title: title, Author: author, Year: year})

	/*err = json.NewEncoder(w).Encode(generate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}*/
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", handleQuote)
	http.HandleFunc("/generate", generateQuote)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
