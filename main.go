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
	// Initialize the SQLite database
	db, err := initDB("url_shortener.db")
	if err != nil {
		fmt.Println("Failed to initialize the database:", err)
		return
	}
	defer db.Close()

	// Parse all templates from the ./templates directory
	tmpl := template.Must(template.New("").ParseGlob("./templates/*"))

	router := http.NewServeMux()

	// Home page route (only handles GET)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleHome(w, r, tmpl)
	})

	// Route for handling the shortening of the URL
	router.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		handleShorten(w, r, tmpl, db)
	})

	// Route for handling short URL redirection (e.g., /abcd1234)
	router.HandleFunc("/{shortKey}", func(w http.ResponseWriter, r *http.Request) {
		handleRedirect(w, r, db)
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Starting website at http://localhost:8080")

	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("An error occurred:", err)
	}
}
