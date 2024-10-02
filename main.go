package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

type PageData struct {
	OriginalUrl string
	ShortUrl    string
	ErrorMsg    string
}

func main() {
	// Parse all templates from the ./templates directory
	tmpl := template.Must(template.New("").ParseGlob("./templates/*"))

	router := http.NewServeMux()

	// Home page route (handles GET and POST for URL submission)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleHome(w, r, tmpl)
	})

	// Route for handling short URL redirection (e.g., /abcd1234)
	router.HandleFunc("/{shortKey}", func(w http.ResponseWriter, r *http.Request) {
		handleRedirect(w, r)
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Starting website at http://localhost:8080")

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("An error occurred:", err)
	}
}
